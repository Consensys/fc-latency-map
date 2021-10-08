/*
Package config - combines operations used to setup parameters for Retrieval Register node in FileCoin network
*/
package config

import (
	"os"

	"github.com/spf13/viper"
)

func NewConfig() *viper.Viper {
	config := viper.New()
	exist := fileExists(".env")
	if exist {
		config.SetConfigFile(".env")
		if err := config.ReadInConfig(); err != nil {
			panic("Failed read config")
		}
	} else {
		config.AutomaticEnv()
	}
	return config
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
