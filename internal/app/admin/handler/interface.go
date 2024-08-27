package handler

import (
	"github.com/gin-gonic/gin"
)

type Store interface {
	Root(ginContext *gin.Context)
	Login(ginContext *gin.Context)
	ProxyInfo(ginContext *gin.Context)
}
