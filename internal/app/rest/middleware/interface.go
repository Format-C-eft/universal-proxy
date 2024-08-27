package middleware

import (
	"github.com/gin-gonic/gin"
)

type Store interface {
	Session() gin.HandlerFunc
	ReadRequestBody() gin.HandlerFunc
	SendAndReceive() gin.HandlerFunc
}
