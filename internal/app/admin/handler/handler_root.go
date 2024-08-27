package handler

import (
	"net/http"

	"github.com/Format-C-eft/utils/headers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *StoreImpl) Root(ginContext *gin.Context) {
	s := sessions.Default(ginContext)
	sessionID := s.Get("sessionID")
	if sessionID == nil {
		ginContext.Writer.Header().Set(headers.Location, "/login")
		ginContext.AbortWithStatus(http.StatusFound)
		return
	}

	ginContext.Writer.Header().Set(headers.Location, "/proxy/info")
	ginContext.AbortWithStatus(http.StatusFound)
}
