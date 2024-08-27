package proxy

import (
	"time"

	"github.com/Format-C-eft/universal-proxy/internal/config"
)

func (s *StoreImpl) GetProxyInfo() []Info {
	results := make([]Info, 0)
	s.proxyMap.Range(func(_, value any) bool {
		valueInfo, _ := value.(*proxy)

		textError := ""
		if valueInfo.LastErr != nil {
			textError = valueInfo.LastErr.Error()
		}

		results = append(results, Info{
			Address:         valueInfo.IP,
			AddressExternal: valueInfo.ExternalIP,
			SessionID:       valueInfo.SessionID,
			Login:           valueInfo.Login,
			LastURL:         valueInfo.LastURL,
			LastUsed:        valueInfo.LastUsed.In(config.DefaultLocation).Format(config.LayoutDate),
			IsActive:        valueInfo.IsActive,
			IsUsedNow:       valueInfo.LastUsed.Add(s.cfg.SessionLifetime).After(time.Now()),
			LastErr:         textError,
			CountErr:        valueInfo.CountErr,
		})

		return true
	})

	return results
}
