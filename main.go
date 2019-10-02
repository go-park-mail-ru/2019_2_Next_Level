package main

import (
	"back/config"
	"back/daemon"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func Configurate(conf *config.Config) error {
	var isLocalhost bool

	flag.BoolVar(&isLocalhost, "local", false, "Is it local mashine")
	flag.Parse()
	if err := (*conf).Inflate(); err != nil {
		log.Println("Cannot read config: ", err)
		return errors.New("Cannot read config")
	}
	if !isLocalhost {
		osPort := os.Getenv("PORT")
		if osPort != "" {
			(*conf).Port = osPort
		}
	}
	if isLocalhost {
		// ex, err := os.Executable()
		// if err != nil {
		// 	panic(err)
		// }
		// location := filepath.Dir(ex)
		(*conf).SelfURL = "http://localhost:" + (*conf).Port
	}
	return nil

}

func main() {
	config.Configuration = config.Config{}
	Configurate(&config.Configuration)

	if err := daemon.Run(&config.Configuration); err != nil {
		fmt.Printf("Error during daemon startup: %s\n", err)
	}
}
