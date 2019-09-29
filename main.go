package main

import (
	"back/daemon"
	"flag"
	"fmt"
	"os"
)

func inflateDaemonConfig() *daemon.Config {
	config := &daemon.Config{}

	flag.StringVar(&config.Port, "port", "80", "Port to listen")
	flag.StringVar(&config.FrontendPath, "front", "./", "Path to frontend to share")
	flag.StringVar(&config.FrontendUrl, "furl", "locahost:3001", "Address of the frontend")
	flag.StringVar(&config.AvatarDirPath, "avadir", "./static/", "Path to the avatars")
	flag.Parse()

	return config
}

func main() {
	config := inflateDaemonConfig()
	config.Port = os.Getenv("PORT")

	if err := daemon.Run(config); err != nil {
		fmt.Printf("Error during daemon startup: %s\n", err)
	}
}
