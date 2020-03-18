package config

import (
	"os"

	"github.com/thomgray/codebook/util"
)

// AppConfig ...
type AppConfig struct {
	HomeDir string
}

var _appConfig *AppConfig = nil

// GetAppConfig ...
func GetAppConfig() AppConfig {
	if _appConfig == nil {
		_appConfig = loadAppConfig()
	}
	return *_appConfig
}

func loadAppConfig() *AppConfig {
	d, err := os.UserHomeDir()
	util.Check(err)

	return &AppConfig{
		HomeDir: d,
	}
}
