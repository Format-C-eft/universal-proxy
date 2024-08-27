package admin

import (
	"github.com/Format-C-eft/universal-proxy/internal/app/admin/handler"
	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/utils/gin/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const (
	headHTML   = "./.fs/html/admin/head.html"
	headerHTML = "./.fs/html/admin/header.html"
	footerHTML = "./.fs/html/admin/footer.html"
	scriptHTML = "./.fs/html/admin/script.html"

	adminTmpl     = "./.fs/html/admin/wallets.tmpl"
	proxyInfoTmpl = "./.fs/html/admin/proxy_info.tmpl"
	walletsTmpl   = "./.fs/html/admin/wallets.tmpl"

	appleTouchIcon  = "./.fs/html/admin/img/apple-touch-icon.png"
	favicon         = "./.fs/html/admin/img/favicon.ico"
	favicon16       = "./.fs/html/admin/img/favicon-16x16.png"
	favicon32       = "./.fs/html/admin/img/favicon-32x32.png"
	safariPinnedTab = "./.fs/html/admin/img/safari-pinned-tab.svg"

	checkGreenIcon = "./.fs/html/admin/img/check-green-50.svg"
	checkRedIcon   = "./.fs/html/admin/img/check-red-50.svg"
)

func New(
	cfg config.AppAdmin,
	handlerStore handler.Store,
) *gin.Engine {
	router := gin.New()

	store := cookie.NewStore([]byte(cfg.SessionPassword))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   1800,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	})

	router.Use(middleware.Recovery())
	router.Use(middleware.ChangeLoggerLevel())
	router.Use(middleware.Logger())
	router.Use(sessions.Sessions(config.CookieSessionAdminName, store))
	router.Use(cors.New(corsConfig))

	router.GET("/", handlerStore.Root)
	router.GET("/proxy/info", handlerStore.ProxyInfo)

	group := router.Group("/", gin.BasicAuth(cfg.Users))
	group.GET("/login", handlerStore.Login)

	router.LoadHTMLFiles(headHTML, headerHTML, footerHTML, scriptHTML, adminTmpl, walletsTmpl, proxyInfoTmpl)

	router.StaticFile("/apple-touch-icon.png", appleTouchIcon)
	router.StaticFile("/favicon.ico", favicon)
	router.StaticFile("/favicon-16x16.png", favicon16)
	router.StaticFile("/favicon-32x32.png", favicon32)
	router.StaticFile("/safari-pinned-tab.svg", safariPinnedTab)

	router.StaticFile("/check-green-50.svg", checkGreenIcon)
	router.StaticFile("/check-red-50.svg", checkRedIcon)

	return router
}
