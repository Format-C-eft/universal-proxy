package proxy

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Format-C-eft/utils/headers"
	"github.com/Format-C-eft/utils/json"
	"github.com/Format-C-eft/utils/logger"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/Format-C-eft/universal-proxy/internal/config"
)

type StoreImpl struct {
	proxyMap *sync.Map
	cfg      *config.ProxyConfig
}

func New(ctx context.Context, cfg config.ProxyConfig) (*StoreImpl, error) {
	storeImpl := &StoreImpl{
		proxyMap: new(sync.Map),
		cfg:      &cfg,
	}

	storeImpl.updatePoolTor(ctx)
	go storeImpl.checkAllProxyActive(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(cfg.UpdateHosts.RefreshTimeout)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				storeImpl.updatePoolTor(ctx)
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(cfg.CheckActive.RefreshTimeout)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Так как эта операция по времени возможно не быстрая,
				// мы таймер сначала останавливаем потом запускаем
				ticker.Stop()

				storeImpl.checkAllProxyActive(ctx)

				// сбросим таймер, что бы не было накладок
				ticker.Reset(cfg.CheckActive.RefreshTimeout)
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(cfg.Restart.RefreshTimeout)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Так как эта операция по времени возможно не быстрая,
				// мы таймер сначала останавливаем потом запускаем
				ticker.Stop()

				storeImpl.restartRemoteProxy()

				// сбросим таймер, что бы не было накладок
				ticker.Reset(cfg.Restart.RefreshTimeout)
			}
		}
	}(ctx)

	return storeImpl, nil
}

func (s *StoreImpl) updatePoolTor(ctx context.Context) {
	ips, err := net.LookupIP(s.cfg.UpdateHosts.Host)
	if err != nil {
		logger.ErrorKV(ctx, "Update pool tor - cant get ip addresses from dns", "err", err)
		return
	}

	if len(ips) == 0 {
		s.proxyMap = new(sync.Map)
		logger.ErrorKV(ctx, "Update pool tor - count TOR ip addresses is zero")
		return
	}

	addressMap := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		address := ip.String()

		addressMap[address] = struct{}{}

		if _, ok := s.proxyMap.Load(address); ok {
			continue
		}

		s.proxyMap.Store(
			address,
			&proxy{
				IP: address,
				Client: resty.New().
					SetProxy(fmt.Sprintf("socks5://%s:%d", address, s.cfg.Client.Port)).
					SetTimeout(s.cfg.Client.ClientTimeout).
					SetRetryCount(1).
					SetRedirectPolicy(
						resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
							return http.ErrUseLastResponse
						}),
					).
					SetDisableWarn(true).
					SetLogger(logger.GetLogger()),
			},
		)
	}

	s.proxyMap.Range(func(key, _ any) bool {
		if _, ok := addressMap[key.(string)]; !ok {
			s.proxyMap.Delete(key)
		}

		return true
	})

	logger.Info(ctx, "Update pool tor - success")
}

func (s *StoreImpl) checkAllProxyActive(ctx context.Context) {
	s.proxyMap.Range(func(key, value any) bool {
		val := value.(*proxy)

		externalIP, err := s.checkProxyActive(val.Client)
		if err != nil {
			// Это если произошла ошибка при проверке, тогда в любом случае помечаем что тут все не очень.
			val.Login = ""
			val.SessionID = ""
			val.LastErr = err
			val.IsActive = false
			val.CountErr++

			s.proxyMap.Store(key, val)

			logger.ErrorKV(ctx, "Check proxy active - error client "+val.IP, "err", err)
			return true
		}

		if !val.IsActive {
			// Это если на прокси висела ошибка, а теперь ее нет
			val.LastErr = nil
			val.IsActive = true
			val.CountErr = 0
			val.ExternalIP = externalIP

			s.proxyMap.Store(key, val)
		}

		// Сюда придет, только если проверка прошла удачно на живой прокси
		return true
	})

	logger.Info(ctx, "Check all proxy active - success")
}

func (s *StoreImpl) checkProxyActive(client *resty.Client) (string, error) {
	req := client.NewRequest()
	req.Header.Set(headers.Accept, "application/json")
	req.Method = http.MethodGet
	req.URL = s.cfg.CheckActive.URL

	resp, err := req.Send()
	if err != nil {
		var val net.Error
		if errors.As(err, &val) && val.Timeout() {
			return "", errors.New("i/o timeout")
		}

		return "", errors.Wrap(err, "req.Send")
	}

	if resp.StatusCode() != http.StatusOK {
		return "", errors.Errorf("StatusCode: %d", resp.StatusCode())
	}

	if !json.GetValues(resp.Body(), "IsTor").ToBool() {
		return "", errors.New("route not in tor network")
	}

	return json.GetValues(resp.Body(), "IP").ToString(), nil
}

func (s *StoreImpl) restartRemoteProxy() {
	s.proxyMap.Range(func(key, value any) bool {
		val, _ := value.(*proxy)
		if val.CountErr >= s.cfg.Restart.ErrorCount {
			if err := s.restartTorFromIP(key.(string)); err != nil {
				val.LastErr = errors.Wrap(err, "restartTorFromIP")
				s.proxyMap.Store(key, val)
			}
		}

		return true
	})
}

func (s *StoreImpl) restartTorFromIP(ip string) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, s.cfg.Restart.Port), s.cfg.Restart.ClientTimeout)
	if err != nil {
		return errors.Wrap(err, "Connecting to tor control port failed")
	}

	_, err = conn.Write([]byte(fmt.Sprintf("authenticate \"%s\"\nsignal newnym\n", s.cfg.Restart.Password)))
	if err != nil {
		return errors.Wrap(err, "Writing command to tor failed")
	}

	torOk := make([]byte, 14)
	_, err = conn.Read(torOk)
	if err != nil {
		return errors.Wrap(err, "Reading tor control port auth response failed")
	}

	if res := bytes.Compare(torOk, []byte("250 OK\r\n250 OK")); res != 0 {
		return errors.Errorf("Tor execute command failed: %s", string(torOk))
	}

	return nil
}
