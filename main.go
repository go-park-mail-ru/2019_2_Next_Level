package main

import (
	"back/config"
	"back/daemon"
	"fmt"
	"log"
)

func main() {
	configuration := config.Config{}
	if err := configuration.Inflate(); err != nil {
		log.Println("Cannot read config: ", err)
		return
	}

	if err := daemon.Run(&configuration); err != nil {
		fmt.Printf("Error during daemon startup: %s\n", err)
	}
}
