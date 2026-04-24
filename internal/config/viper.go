package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	config.AddConfigPath("../")
	config.AddConfigPath("../../")

	err := config.ReadInConfig()
	if err != nil {
		fmt.Printf("Warning: No config file found, using defaults and environment variables: %v\n", err)
	}

	config.AutomaticEnv()

	return config
}

func NewTestViper() *viper.Viper {
	config := viper.New()
	config.SetConfigName("config_test")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	config.AddConfigPath("../")
	config.AddConfigPath("../../")

	err := config.ReadInConfig()
	if err != nil {
		fmt.Printf("Warning: No test config file found, using defaults and environment variables: %v\n", err)
	}

	config.AutomaticEnv()

	return config
}
