package status

import (
	"net/http"
	"sync/atomic"

	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/gin-gonic/gin"
)

func New(isReady *atomic.Value) *gin.Engine {
	h := &handlers{
		isReady: isReady,
	}

	router := gin.New()

	router.GET("/live", h.live)
	router.GET("/ready", h.ready)
	router.GET("/version", h.version)

	return router
}

type handlers struct {
	isReady *atomic.Value
}

func (h *handlers) live(ginContext *gin.Context) {
	ginContext.Status(http.StatusOK)
}

func (h *handlers) ready(ginContext *gin.Context) {
	if h.isReady == nil || !h.isReady.Load().(bool) {
		ginContext.Status(http.StatusServiceUnavailable)
		return
	}

	ginContext.Status(http.StatusOK)
}

func (h *handlers) version(ginContext *gin.Context) {
	ginContext.JSON(
		http.StatusOK,
		config.GetVersion(),
	)
}
