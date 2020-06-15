package main

import (
	"os"
	"path/filepath"

	"github.com/thomgray/codebook/config"
)

func install() {
	installDir := filepath.Join(config.GetAppConfig().HomeDir, ".codebook")

	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		os.Mkdir(installDir, os.ModePerm)
	}
}
