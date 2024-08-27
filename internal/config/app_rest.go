package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	AppAddress struct {
		Scheme string
		Host   string
	}

	AppRest struct {
	}

	AppRestMiddleware struct {
		CurrentAddress        AppAddress
		StartAddress          AppAddress
		SessionCookieDuration time.Duration
		SessionCache          LocalCache
	}

	AppRestHandler struct {
	}
)

func GetAppRestMiddleware() AppRestMiddleware {
	return AppRestMiddleware{
		CurrentAddress: AppAddress{
			Scheme: viper.GetString(AppRestMiddlewareCurrentAddressScheme),
			Host:   viper.GetString(AppRestMiddlewareCurrentAddressHost),
		},
		StartAddress: AppAddress{
			Scheme: viper.GetString(AppRestMiddlewareStartAddressScheme),
			Host:   viper.GetString(AppRestMiddlewareStartAddressHost),
		},
		SessionCookieDuration: viper.GetDuration(AppRestMiddlewareSessionCookieDuration),
		SessionCache: LocalCache{
			Name: viper.GetString(AppRestMiddlewareSessionLocalCacheName),
			Size: viper.GetInt(AppRestMiddlewareSessionLocalCacheSize),
			TTL:  viper.GetDuration(AppRestMiddlewareSessionLocalCacheTTL),
		},
	}
}

func GetAppRestHandler() AppRestHandler {
	return AppRestHandler{}
}

func GetAppRest() AppRest {
	return AppRest{}
}
