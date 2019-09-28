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
	flag.StringVar(&config.FrontendUrl, "furl", "locahost:3001", "Address of the frontend")
	flag.StringVar(&config.AvatarDirPath, "avadir", "./static/avatar", "Path to the avatars")
	flag.Parse()

	return config
}

func main() {
	config := inflateDaemonConfig()

	if err := daemon.Run(config); err != nil {
		fmt.Printf("Error during daemon startup: %s\n", err)
	}
}
