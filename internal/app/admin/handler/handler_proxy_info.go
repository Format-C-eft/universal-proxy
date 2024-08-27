package handler

import (
	"bytes"
	"net"
	"net/http"
	"sort"

	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/utils/headers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *StoreImpl) ProxyInfo(ginContext *gin.Context) {
	s := sessions.Default(ginContext)
	sessionID := s.Get("sessionID")
	if sessionID == nil {
		ginContext.Writer.Header().Set(headers.Location, "/login")
		ginContext.AbortWithStatus(http.StatusFound)
		return
	}

	proxyInfo := h.proxyStore.GetProxyInfo()

	sort.Slice(proxyInfo, func(i, j int) bool {
		return bytes.Compare(net.ParseIP(proxyInfo[i].Address), net.ParseIP(proxyInfo[j].Address)) < 0
	})

	ginContext.HTML(http.StatusOK, "proxy_info.tmpl", gin.H{
		"info":     config.GetVersion(),
		"sessions": proxyInfo,
	})
}
