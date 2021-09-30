package config

import (
	"github.com/spf13/viper"
)

// NewMockConfig creates a new mock instance.
func NewMockConfig() *viper.Viper {
	mock := viper.New()
	return mock
}
