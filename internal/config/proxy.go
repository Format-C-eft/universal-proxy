package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	ProxyConfig struct {
		SessionLifetime time.Duration

		UpdateHosts ProxyUpdateHosts
		CheckActive ProxyCheckActive
		Client      ProxyClient
		Restart     ProxyRestart
	}

	ProxyUpdateHosts struct {
		Host           string
		RefreshTimeout time.Duration
	}

	ProxyCheckActive struct {
		URL            string
		RefreshTimeout time.Duration
		ClientTimeout  time.Duration
	}

	ProxyRestart struct {
		Port           int64
		Password       string
		ErrorCount     int
		RefreshTimeout time.Duration
		ClientTimeout  time.Duration
	}

	ProxyClient struct {
		Port             int64
		ConnectTimeout   time.Duration
		ConnectKeepAlive time.Duration
		ClientTimeout    time.Duration
	}
)

func GetProxyConfig() ProxyConfig {
	return ProxyConfig{
		SessionLifetime: viper.GetDuration(ProxySessionLifetime),
		UpdateHosts: ProxyUpdateHosts{
			Host:           viper.GetString(ProxyUpdateHostsHost),
			RefreshTimeout: viper.GetDuration(ProxyUpdateHostsRefreshTimeout),
		},
		CheckActive: ProxyCheckActive{
			URL:            viper.GetString(ProxyCheckActiveURL),
			RefreshTimeout: viper.GetDuration(ProxyCheckActiveRefreshTimeout),
			ClientTimeout:  viper.GetDuration(ProxyCheckActiveClientTimeout),
		},
		Client: ProxyClient{
			Port:             viper.GetInt64(ProxyClientPort),
			ConnectTimeout:   viper.GetDuration(ProxyClientConnectTimeout),
			ConnectKeepAlive: viper.GetDuration(ProxyClientConnectKeepAlive),
			ClientTimeout:    viper.GetDuration(ProxyClientClientTimeout),
		},
		Restart: ProxyRestart{
			Port:           viper.GetInt64(ProxyRestartPort),
			Password:       viper.GetString(ProxyRestartPassword),
			ErrorCount:     viper.GetInt(ProxyRestartErrorCount),
			RefreshTimeout: viper.GetDuration(ProxyRestartRefreshTimeout),
			ClientTimeout:  viper.GetDuration(ProxyRestartClientTimeout),
		},
	}
}
