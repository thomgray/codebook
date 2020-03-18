package main

import (
	"log"
	"os"

	"github.com/thomgray/codebook/config"
	"github.com/thomgray/codebook/controller"
)

// var config *config.Config

func main() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	install()
	config := config.MakeConfig()
	controller := controller.InitMainController(config)
	defer file.Close()
	defer controller.Start()
}
