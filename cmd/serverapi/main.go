package main

import (
	serverapiconfig "2019_2_Next_Level/internal/serverapi/config"
	"2019_2_Next_Level/internal/serverapi/log"
	"2019_2_Next_Level/internal/serverapi/server"
	"2019_2_Next_Level/pkg/config"
	"2019_2_Next_Level/pkg/logger"
	"flag"
	"sync"
)

const (
	configFilenameDefault = "http_service.config.json"
)

func main() {
	log.SetLogger(logger.NewLog())
	log.Log().SetPrefix("HttpService")
	log.Log().I("API Server started. Hello!")

	if err := initializeConfig(); err != nil {
		log.Log().E(err)
		return
	}
	// curl -d "to=andrey" http://localhost:3001/mail/send

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go server.Run(wg)
	wg.Wait()

}

func initializeConfig() error {
	configFilename := flag.String("config", configFilenameDefault, "Path to config file")
	dbUser := flag.String("dbuser", "", "User for database")
	dbPassword := flag.String("dbpass", "", "Password for database")
	flag.Parse()

	return config.Inflate(*configFilename, &serverapiconfig.Conf, *dbUser, *dbPassword)
}
