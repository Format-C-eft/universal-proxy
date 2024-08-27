package config

import (
	"context"
	"flag"
	"strings"

	"github.com/Format-C-eft/utils/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	LocalRunKey = "local-run"
	VersionKey  = "version"
)

var (
	LocalRunFlag = false
	VersionFlag  = false
)

func init() {
	viper.SetEnvPrefix(strings.ReplaceAll(AppName, "-", "_"))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")

	flag.BoolVar(&LocalRunFlag, LocalRunKey, LocalRunFlag, "enable development config YAML parser")
	flag.BoolVar(&VersionFlag, VersionKey, VersionFlag, "show version")
	flag.Parse()

	if LocalRunFlag {
		viper.SetConfigName("values_local")
	} else {
		viper.SetConfigName("values")
	}

	viper.AddConfigPath(".fs/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.FatalKV(context.Background(), "cant read config", "err", err)
	}

	logLevel := viper.GetString(LogLevel)
	if zapLevel, ok := logger.StrToZapLevel(logLevel); ok {
		logger.SetLevel(zapLevel)
		logger.InfoF(context.Background(), "Logger set level - %s", logLevel)
	}

	viper.WatchConfig()

	viper.OnConfigChange(
		func(_ fsnotify.Event) {
			watchLogLevel := viper.GetString(LogLevel)
			if zapLevel, ok := logger.StrToZapLevel(watchLogLevel); ok {
				logger.SetLevel(zapLevel)
				logger.InfoF(context.Background(), "Logger set level - %s", watchLogLevel)
			}
		},
	)
}

func GetPort(portName string) int64 {
	return viper.GetInt64("PORT_" + portName)
}

var (
	LogLevel = "LOG_LEVEL"

	ProxySessionLifetime           = "PROXY_SESSION_LIFETIME"
	ProxyUpdateHostsHost           = "PROXY_UPDATE_HOSTS_HOST"
	ProxyUpdateHostsRefreshTimeout = "PROXY_UPDATE_HOSTS_REFRESH_TIMEOUT"
	ProxyClientPort                = "PROXY_CLIENT_PORT"
	ProxyClientConnectTimeout      = "PROXY_CLIENT_CONNECT_TIMEOUT"
	ProxyClientConnectKeepAlive    = "PROXY_CLIENT_CONNECT_KEEPALIVE"
	ProxyClientClientTimeout       = "PROXY_CLIENT_TIMEOUT"
	ProxyRestartPort               = "PROXY_RESTART_PORT"
	ProxyRestartPassword           = "PROXY_RESTART_PASSWORD"
	ProxyRestartErrorCount         = "PROXY_RESTART_ERROR_COUNT"
	ProxyRestartRefreshTimeout     = "PROXY_RESTART_REFRESH_TIMEOUT"
	ProxyRestartClientTimeout      = "PROXY_RESTART_CLIENT_TIMEOUT"
	ProxyCheckActiveURL            = "PROXY_CHECK_ACTIVE_URL"
	ProxyCheckActiveRefreshTimeout = "PROXY_CHECK_ACTIVE_REFRESH_TIMEOUT"
	ProxyCheckActiveClientTimeout  = "PROXY_CHECK_ACTIVE_CLIENT_TIMEOUT"

	AppRestMiddlewareStartAddressScheme   = "APP_REST_MIDDLEWARE_START_ADDRESS_SCHEME"
	AppRestMiddlewareStartAddressHost     = "APP_REST_MIDDLEWARE_START_ADDRESS_HOST"
	AppRestMiddlewareCurrentAddressScheme = "APP_REST_MIDDLEWARE_CURRENT_ADDRESS_SCHEME"
	AppRestMiddlewareCurrentAddressHost   = "APP_REST_MIDDLEWARE_CURRENT_ADDRESS_HOST"

	AppRestMiddlewareSessionCookieDuration = "APP_REST_MIDDLEWARE_SESSION_COOKIE_DURATION"

	AppRestMiddlewareSessionLocalCacheName = "APP_REST_MIDDLEWARE_SESSION_LOCAL_CACHE_NAME"
	AppRestMiddlewareSessionLocalCacheSize = "APP_REST_MIDDLEWARE_SESSION_LOCAL_CACHE_SIZE"
	AppRestMiddlewareSessionLocalCacheTTL  = "APP_REST_MIDDLEWARE_SESSION_LOCAL_CACHE_TTL"

	AppAdminUsers = "APP_ADMIN_USERS"

	AppAdminSessionPassword = "APP_ADMIN_SESSION_PASSWORD" //nolint
)
