package middleware

import (
	"math"
	"net/http"

	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type sessionInfo struct {
	UserUUID           uuid.UUID
	UserLogin          string
	RequestPreviousURI config.AppAddress
	RequestURI         config.AppAddress
	CookiesMap         map[string]map[string]http.Cookie
}

func (m *StoreImpl) Session() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		sessionID, err := ginContext.Cookie(config.CookieSessionName)
		if err != nil {
			logger.WarnKV(ginContext, "Session: Request no cookie sessionID, start new session")

			u, _ := uuid.NewUUID()
			sessionID = u.String()
		}

		sessionRaw, err := m.sessionCache.Get(ginContext, sessionID)
		if err != nil {
			logger.ErrorKV(ginContext, "Session: sessionCache.Get", "sessionID", sessionID, "err", err)
		}

		session := &sessionInfo{
			CookiesMap: make(map[string]map[string]http.Cookie, 1),
		}

		if sessionCast, ok := sessionRaw.(*sessionInfo); ok {
			session = sessionCast
		} else {
			logger.ErrorKV(ginContext, "Session: cant cast cache value to sessionInfo", "cacheValue", sessionRaw)
		}

		if session.RequestURI.Host == "" {
			session.RequestURI = m.cfg.StartAddress
		}

		// Создаем мапку с куками для хоста запроса
		if _, ok := session.CookiesMap[session.RequestURI.Host]; !ok {
			session.CookiesMap[session.RequestURI.Host] = make(map[string]http.Cookie, 10)
		}

		// Устанавливаем куку сессии
		ginContext.SetCookie(
			config.CookieSessionName,
			sessionID,
			int(math.Round(m.cfg.SessionCookieDuration.Seconds())),
			"/",
			"",
			false,
			false,
		)

		// Докидываем в контекст выполнения запроса все необходимые для этого параметры зависящие от сессии
		ginContext.Set(config.ContextSessionKey, sessionID)
		ginContext.Set(config.ContextUserUUIDKey, session.UserUUID)
		ginContext.Set(config.ContextUserLoginKey, session.UserLogin)
		ginContext.Set(config.ContextRequestPreviousURIKey, session.RequestPreviousURI)
		ginContext.Set(config.ContextRequestURIKey, session.RequestURI)
		ginContext.Set(config.ContextCookiesMapKey, session.CookiesMap[session.RequestURI.Host])
		ginContext.Set(config.ContextSessionIsNeedDelete, false)

		ginContext.Next()

		// Нужно только для одного кейса, когда надо выкинуть пользователя с концами.
		// Для этого просто надо его сессию прибить у нас. У пользователя ничего нет кроме, куки с нашей сессией
		if ginContext.MustGet(config.ContextSessionIsNeedDelete).(bool) {
			if errDel := m.sessionCache.Remove(ginContext, sessionID); errDel != nil {
				logger.ErrorKV(ginContext, "Session: sessionCache.Remove", "err", errDel)
			}
			return
		}

		// Запоминаем данные из контекста выполнения.
		session.UserUUID = ginContext.MustGet(config.ContextUserUUIDKey).(uuid.UUID)
		session.UserLogin = ginContext.MustGet(config.ContextUserLoginKey).(string)
		session.RequestPreviousURI = ginContext.MustGet(config.ContextRequestPreviousURIKey).(config.AppAddress)
		session.RequestURI = ginContext.MustGet(config.ContextRequestURIKey).(config.AppAddress)

		// Сохраняем в кеш данные в любом случае даже если они не изменились, что бы обновить TTL в кеше
		if errSet := m.sessionCache.Set(ginContext, sessionID, session); errSet != nil {
			logger.ErrorKV(ginContext, "Session: sessionCache.Set", "err", errSet)
		}
	}
}
