package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/universal-proxy/internal/proxy"
	"github.com/Format-C-eft/utils/headers"
	"github.com/Format-C-eft/utils/logger"
	"github.com/gin-gonic/gin"
)

var (
	regexpPatternCookieBodyV1   = regexp.MustCompile(`document\.cookie\s*=\s*\"(.*)=".*expires\s*=\s*(.*?);\s*path\s*=\s*(.*?);`)
	regexpPatternLocationBodyV1 = regexp.MustCompile(`location\.href\s*=\s*\"((http|https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-]))\"`)
)

func (m *StoreImpl) SendAndReceive() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		requestPreviousURI := ginContext.MustGet(config.ContextRequestPreviousURIKey).(config.AppAddress)
		requestURI := ginContext.MustGet(config.ContextRequestURIKey).(config.AppAddress)
		cookiesMap := ginContext.MustGet(config.ContextCookiesMapKey).(map[string]http.Cookie)
		requestBody := ginContext.MustGet(config.ContextRequestBodyKey).([]byte)

		req := &proxy.Request{
			Method:     ginContext.Request.Method,
			RequestURI: fmt.Sprintf("%s://%s%s", requestURI.Scheme, requestURI.Host, ginContext.Request.RequestURI),
			Headers:    fixAndGetRequestHeaders(ginContext.Request.Header, requestURI, requestPreviousURI),
			Cookies:    fixAndGetRequestCookies(cookiesMap),
			Body:       requestBody,
			User: proxy.User{
				SessionID: ginContext.MustGet(config.ContextSessionKey).(string),
				Login:     ginContext.MustGet(config.ContextUserLoginKey).(string),
			},
		}

		resp, err := m.proxyStore.ExecuteRequest(req)
		if err != nil {
			logger.ErrorKV(ginContext, "SendAndReceive: proxyStore.ExecuteRequest error", "err", err)
			ginContext.HTML(http.StatusInternalServerError, "error.html", gin.H{})
			ginContext.Abort()
			return
		}

		// Сохраняем у нас значение кук из заголовков
		for _, cookie := range resp.Cookies {
			cookiesMap[cookie.Name] = cookie
		}

		// Сохраняем у нас значения кук из тела ответа
		for _, cookie := range getCookieFromBody(resp.Body) {
			cookiesMap[cookie.Name] = cookie
		}

		// Сохраняем URI куда сейчас ходили, для исправления заголовка Referer
		ginContext.Set(config.ContextRequestPreviousURIKey, requestURI)

		// Получаем и исправляем редирект в заголовках
		if valueRaw := resp.Headers.Get(headers.Location); valueRaw != "" {
			value, _ := url.Parse(valueRaw)
			// Если нет хоста, то это редирект на страницу внутри хоста, дальнейшая обработка не уместна
			if value.Host != "" {
				// Нам из всего адреса нужна только схема и хост, с остальным к нам придут
				requestURI.Scheme = value.Scheme
				requestURI.Host = value.Host

				// Подменяем значения на наши
				value.Scheme = m.cfg.CurrentAddress.Scheme
				value.Host = m.cfg.CurrentAddress.Host

				// Перезаписываем заголовок
				resp.Headers.Set(headers.Location, value.String())
			}
		}

		// Получаем и исправляем редирект в теле ответа
		if match := regexpPatternLocationBodyV1.FindStringSubmatch(string(resp.Body)); len(match) != 0 {
			requestURI.Scheme = match[2]
			requestURI.Host = match[3]

			resp.Body = bytes.ReplaceAll(resp.Body, []byte(match[2]+"://"+match[3]), []byte(m.cfg.CurrentAddress.Scheme+"://"+m.cfg.CurrentAddress.Host))
		}

		ginContext.Set(config.ContextRequestURIKey, requestURI)
		ginContext.Set(config.ContextRequestKey, req)
		ginContext.Set(config.ContextResponseKey, resp)
		ginContext.Set(config.ContextReturnErrorKey, false)

		ginContext.Next()

		if ginContext.MustGet(config.ContextReturnErrorKey).(bool) {
			ginContext.HTML(http.StatusInternalServerError, "error.html", gin.H{})
			ginContext.Abort()
			return
		}

		// Удаляем заголовки
		resp.Headers.Del(headers.TransferEncoding) // gin сильно ругается на это
		resp.Headers.Del(headers.ContentEncoding)  // Тут лежит значение gzip, а мы и так распаковываем

		// Устанавливаем заголовки
		for key, values := range resp.Headers {
			for _, header := range values {
				ginContext.Header(key, header)
			}
		}

		// Устанавливаем контент тип и размер
		ginContext.Header(headers.ContentLength, strconv.Itoa(len(resp.Body)))

		// Отвечаем клиенту
		ginContext.String(resp.StatusCode, string(resp.Body))
	}
}

func fixAndGetRequestHeaders(header http.Header, requestAddress, requestPreviousAddress config.AppAddress) http.Header {
	if requestPreviousAddress.Scheme == "" {
		requestPreviousAddress = requestAddress
	}

	result := make(http.Header, len(header))
	for key, value := range header {
		switch key {
		case headers.Cookie:
			continue
		case headers.Referer:
			if requestAddress.Host == requestPreviousAddress.Host {
				refererURL, _ := url.Parse(value[0])
				refererURL.Scheme = requestPreviousAddress.Scheme
				refererURL.Host = requestPreviousAddress.Host

				result.Set(headers.Referer, refererURL.String())

				continue
			}

			result.Set(headers.Referer, fmt.Sprintf("%s://%s/", requestPreviousAddress.Scheme, requestPreviousAddress.Host))
		case headers.Origin:
			result.Set(headers.Origin, requestAddress.Scheme+"://"+requestAddress.Host)
		case headers.Host:
			result.Set(headers.Host, requestAddress.Host)
		default:
			for _, s := range value {
				result.Add(key, s)
			}
		}
	}

	return result
}

func fixAndGetRequestCookies(innerCookiesMap map[string]http.Cookie) []*http.Cookie {
	// Сделано именно тут, для случая если потребуется добавлять куки других адресов
	result := make([]*http.Cookie, 0, len(innerCookiesMap))

	for _, cookie := range innerCookiesMap {
		// Если кука просрочена, то мы не отправляем её
		if !cookie.Expires.IsZero() && cookie.Expires.In(time.UTC).Before(time.Now()) {
			continue
		}

		result = append(result,
			&http.Cookie{
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

func getCookieFromBody(body []byte) []http.Cookie {
	result := make([]http.Cookie, 0, 1)
	if len(body) == 0 {
		return result
	}

	for _, values := range regexpPatternCookieBodyV1.FindAllStringSubmatch(string(body), -1) {
		if len(values) != 4 {
			continue
		}

		expires, _ := time.Parse("Mon, 02-Jan-06 15:04:05 MST", values[2])

		result = append(result, http.Cookie{
			Name:       values[1],
			Path:       strings.ReplaceAll(values[3], "\"", ""),
			Expires:    expires,
			RawExpires: values[2],
		})
	}

	return result
}
