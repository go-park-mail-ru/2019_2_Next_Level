package main

import (
	"flag"
	"fmt"

	"./daemon"
)

func inflateDaemonConfig() *daemon.Config {
	config := &daemon.Config{}

	flag.IntVar(&config.Port, "port", 3000, "Port to listen")
	flag.StringVar(&config.FrontendPath, "front", "./", "Path to frontend to share")
	flag.Parse()

	return config
}

func main() {
	config := inflateDaemonConfig()

	if err := daemon.Run(config); err != nil {
		fmt.Printf("Error during daemon startup: %s\n", err)
	}
}
