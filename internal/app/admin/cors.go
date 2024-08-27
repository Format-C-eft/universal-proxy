package admin

import (
	"time"

	"github.com/gin-contrib/cors"
)

var corsConfig cors.Config

func init() {
	corsConfig = cors.Config{
		AllowAllOrigins:        true,
		AllowOrigins:           nil,
		AllowOriginFunc:        nil,
		AllowMethods:           []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTION"},
		AllowHeaders:           []string{"Accept", "Content-Type", "Content-Length", "Origin", "Authorization"},
		AllowCredentials:       true,
		ExposeHeaders:          nil,
		MaxAge:                 600 * time.Second,
		AllowWildcard:          false,
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             false,
	}
}
