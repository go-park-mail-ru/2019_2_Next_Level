package main

import (
	"2019_2_Next_Level/internal/support"
	supportconfig "2019_2_Next_Level/internal/support/config"
	"2019_2_Next_Level/internal/support/log"
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
	log.Log().SetPrefix("SupportService")
	log.Log().I("Server started. Hello!")

	//if err := initializeConfig(); err != nil {
	//	log.Log().E(err)
	//	return
	//}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go support.Run(wg)
	wg.Wait()

}

func initializeConfig() error {
	configFilename := flag.String("config", configFilenameDefault, "Path to config file")
	dbUser := flag.String("dbuser", "", "User for database")
	dbPassword := flag.String("dbpass", "", "Password for database")
	flag.Parse()

	return config.Inflate(*configFilename, &supportconfig.Conf, *dbUser, *dbPassword)
}