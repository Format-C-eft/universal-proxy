package rest

import (
	"github.com/Format-C-eft/utils/gin/middleware"
	"github.com/gin-gonic/gin"

	appRestHandlers "github.com/Format-C-eft/universal-proxy/internal/app/rest/handler"
	appRestMiddleware "github.com/Format-C-eft/universal-proxy/internal/app/rest/middleware"
	"github.com/Format-C-eft/universal-proxy/internal/config"
)

const (
	htmlError = "./.fs/html/rest/error.html"
)

func New(middlewareStore appRestMiddleware.Store, _ appRestHandlers.Store, _ config.AppRest) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(middleware.Recovery())
	router.Use(middleware.ChangeLoggerLevel())
	router.Use(middleware.Logger())

	router.Use(middlewareStore.Session())         // Обрабатываем информацию о сессиях
	router.Use(middlewareStore.ReadRequestBody()) // Сделано для того, что бы тело запроса сохранялось в запросе
	router.Use(middlewareStore.SendAndReceive())  // Отправляем запросы и отдаем ответы клиенту только отсюда

	router.LoadHTMLFiles(htmlError)

	return router
}
