package main

import (
	"2019_2_Next_Level/internal/serverapi"
	"2019_2_Next_Level/internal/serverapi/server"
	"2019_2_Next_Level/pkg/config"
	"flag"
	"fmt"
	"log"
	"sync"
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
	// var a post.Sender
	// a = &serverapi.QueueClient{}
	// serverapi.SetQueue(a)
	// serverapi.Run()

	// go server.Run()

	// curl -d "to=andrey" http://localhost:3001/mail/send

	incomingMailHandler := incommail.Secretary{}
	incomingMailHandler.Init()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go server.Run(wg)
	go incomingMailHandler.Run(wg)
	wg.Wait()

	// if err := daemon.Run(&config.Configuration); err != nil {
	// 	fmt.Printf("Error during daemon startup: %s\n", err)
	// }
}

func initializeConfig() error {
	configFilename := flag.String("config", configFilenameDefault, "Path to config file")
	flag.Parse()

	return config.Configurator.Inflate(*configFilename, &serverapi.Conf)
}
