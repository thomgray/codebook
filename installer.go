package main

import (
	"fmt"
	"os"

	"github.com/thomgray/codebook/config"
)

func install() {
	installDir := fmt.Sprintf("%s/.codebook", config.GetAppConfig().HomeDir)

	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		os.Mkdir(installDir, os.ModePerm)
	}
}
