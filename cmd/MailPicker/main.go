package main

import (
	mailpicker "2019_2_Next_Level/internal/MailPicker"
	localconfig "2019_2_Next_Level/internal/MailPicker/config"
	"2019_2_Next_Level/internal/MailPicker/repository"
	postinterface "2019_2_Next_Level/internal/postInterface"
	"2019_2_Next_Level/pkg/config"
	"flag"
	"fmt"
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

	postgresRepo := repository.NewPostgresRepository()
	if postgresRepo == nil {
		fmt.Println("Error during init repo")
		return
	}
	smtpInterface := postinterface.NewQueueClient(localconfig.Conf.RemoteHost, localconfig.Conf.RemotePort)
	quitChan := make(chan interface{}, 1)

	module := mailpicker.NewSecretary(smtpInterface, postgresRepo, quitChan).Init()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go module.Run(wg)
	wg.Wait()
}

func initializeConfig() error {
	configFilename := flag.String("config", configFilenameDefault, "Path to config file")
	dbUser := flag.String("dbuser", "", "User for database")
	dbPassword := flag.String("dbpass", "", "Password for database")
	flag.Parse()

	return config.Inflate(*configFilename, &localconfig.Conf, *dbUser, *dbPassword)
}
