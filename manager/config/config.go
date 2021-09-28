/*
Package config - combines operations used to setup parameters for Retrieval Register node in FileCoin network
*/
package config

import (
	"github.com/spf13/viper"
)

func NewConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigFile(".env")
	if err := config.ReadInConfig(); err != nil {
		panic("Failed read config")
	}

	return config
}
