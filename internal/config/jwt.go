package config

import "github.com/spf13/viper"

// SecretKey reads the JWT secret key from viper config
func SecretKey(v *viper.Viper) string {
	return v.GetString("secretkey")
}
