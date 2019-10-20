package main

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/serverapi"
	"2019_2_Next_Level/pkg/config"
	"flag"
	"fmt"
	"log"
)

const (
	configFilenameDefault = "http_service.config.json"
)

func main() {
	fmt.Println("API Server started. Hello!")
	if err := initializeConfig(); err != nil {
		log.Println(err)
		return
	}
	var a post.Sender
	a = &serverapi.QueueClient{}
	serverapi.SetQueue(a)
	serverapi.Run()

	// if err := daemon.Run(&config.Configuration); err != nil {
	// 	fmt.Printf("Error during daemon startup: %s\n", err)
	// }
}

func initializeConfig() error {
	configFilename := flag.String("config", configFilenameDefault, "Path to config file")
	flag.Parse()

	return config.Configurator.Inflate(*configFilename, &serverapi.Conf)
}
