package config

import "github.com/spf13/viper"

func AppPort(v *viper.Viper) int {
	return v.GetInt("app.port")
}
