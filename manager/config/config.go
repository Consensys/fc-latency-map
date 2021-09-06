/*
Package config - combines operations used to setup parameters for Retrieval Register node in FileCoin network
*/
package config

import (
	"github.com/spf13/viper"
)

func Config() *viper.Viper {
	config := viper.New()
	config.SetConfigFile(".env")
	config.ReadInConfig()
	return config
}
