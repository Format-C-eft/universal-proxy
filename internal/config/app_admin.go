package config

import (
	"github.com/spf13/viper"
)

type AppAdmin struct {
	Users           map[string]string
	SessionPassword string
}

func GetAppAdmin() AppAdmin {
	return AppAdmin{
		Users:           viper.GetStringMapString(AppAdminUsers),
		SessionPassword: viper.GetString(AppAdminSessionPassword),
	}
}
