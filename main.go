package main

import (
	"log"
	"os"

	"github.com/thomgray/codebook/config"
	"github.com/thomgray/codebook/controller"
)

// var config *config.Config

func main() {
	os.Setenv("TCELL_TRUECOLOR", "disable")
	devMode := os.Getenv("codebookdevmode")
	if devMode == "true" {
		file, err := os.OpenFile("info.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
		defer file.Close()
	} else {
		output, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
		log.SetOutput(output)
		defer output.Close()
	}
	install()
	config := config.MakeConfig()
	controller := controller.InitMainController(config)
	defer controller.Start()
}
