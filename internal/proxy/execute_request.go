package proxy

import (
	"net/http"
	"time"

	"github.com/Format-C-eft/utils/headers"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var ErrNoFreeClient = errors.New("no free client available")

func (s *StoreImpl) ExecuteRequest(request *Request) (*Response, error) {
	freeClient, err := s.getFreeClient(request.User, request.RequestURI)
	if err != nil {
		return nil, err
	}

	req := freeClient.Client.NewRequest()

	req.Method = request.Method
	req.URL = request.RequestURI
	req.SetBody(request.Body)

	for _, cookie := range request.Cookies {
		req.SetCookie(cookie)
	}

	for key, values := range request.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, errSend := req.Send()
	if errSend != nil {
		errSend = errors.Wrap(errSend, "proxy.Client.DoTimeout")

		freeClient.IsActive = false
		freeClient.LastErr = errSend
		freeClient.CountErr++

		s.proxyMap.Store(freeClient.IP, freeClient)

		return nil, errSend
	}

	return convertRestyResponseToResponse(resp), nil
}

func (s *StoreImpl) getFreeClient(user User, url string) (*proxy, error) {
	var (
		freeProxyInfoKey string
		freeProxyInfo    *proxy
	)

	s.proxyMap.Range(func(key, value any) bool {
		val, _ := value.(*proxy)
		if !val.IsActive {
			return true
		}

		if freeProxyInfoKey == "" && val.LastUsed.Add(s.cfg.SessionLifetime).Before(time.Now()) {
			freeProxyInfoKey = key.(string)
			freeProxyInfo = val
		}

		if val.SessionID == user.SessionID {
			freeProxyInfoKey = key.(string)
			freeProxyInfo = val
			return false
		}

		return true
	})

	if freeProxyInfoKey == "" {
		return nil, ErrNoFreeClient
	}

	freeProxyInfo.SessionID = user.SessionID
	freeProxyInfo.Login = user.Login
	freeProxyInfo.LastURL = url
	freeProxyInfo.LastUsed = time.Now()

	s.proxyMap.Store(freeProxyInfoKey, freeProxyInfo)

	return freeProxyInfo, nil
}

func convertRestyResponseToResponse(response *resty.Response) *Response {
	// Получаем тело ответа
	result := &Response{
		StatusCode: response.StatusCode(),
		Body:       response.Body(),
		Headers:    make(http.Header),
	}

	// Переливаем все заголовки кроме, контент типа (*не помню почему)
	for key, values := range response.Header() {
		if key == headers.SetCookie {
			continue
		}

		for _, value := range values {
			result.Headers.Add(key, value)
		}
	}

	// Переливаем что бы не было коллизии
	responseCookies := response.Cookies()
	result.Cookies = make([]http.Cookie, len(responseCookies))
	for _, cookie := range responseCookies {
		result.Cookies = append(result.Cookies,
			http.Cookie{
				Name:       cookie.Name,
				Value:      cookie.Value,
				Path:       cookie.Path,
				Domain:     cookie.Domain,
				Expires:    cookie.Expires,
				RawExpires: cookie.RawExpires,
				MaxAge:     cookie.MaxAge,
				Secure:     cookie.Secure,
				HttpOnly:   cookie.HttpOnly,
				SameSite:   cookie.SameSite,
				Raw:        cookie.Raw,
				Unparsed:   cookie.Unparsed,
			},
		)
	}

	return result
}
