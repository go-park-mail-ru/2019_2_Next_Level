package main

import (
	mailpicker "2019_2_Next_Level/internal/MailPicker"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/pkg/config"
	"flag"
	"log"
	"sync"
)

const (
	configFilenameDefault = "mailpicker.config.json"
)

func main() {
	err := initializeConfig()
	if err != nil {
		log.Println(err)
		return
	}

	dbConnection := model.Connection{}
	// repo := repository.NewRepository(&dbConnection)
	// qInterface := postinterface.QueueClient{
	// 	RemoteHost: mailpicker.Conf.RemoteHost,
	// 	RemotePort: mailpicker.Conf.RemotePort,
	// }
	// task := mailpicker.Secretary{}
	// task.Init(&repo, &qInterface)
	// task.DefaultInit(&dbConnection)
	task := mailpicker.NewInstanceDefault(&dbConnection)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go task.Run(wg)
	wg.Wait()
}

func initializeConfig() error {
	configFilename := flag.String("config", configFilenameDefault, "Path to config file")
	flag.Parse()

	return config.Configurator.Inflate(*configFilename, &mailpicker.Conf)
}
