package main

import (
	"back/daemon"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func inflateDaemonConfig() *daemon.Config {
	config := &daemon.Config{}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location := filepath.Dir(ex)

	flag.StringVar(&config.Port, "port", "80", "Port to listen")
	flag.StringVar(&config.FrontendPath, "front", "./", "Path to frontend to share")
	flag.StringVar(&config.FrontendUrl, "furl", "locahost:3001", "Address of the frontend")
	flag.StringVar(&config.AvatarDirPath, "avadir", "avatar/", "Path to the avatars")
	flag.StringVar(&config.StaticDirPath, "staticdir", location+"/public", "Path to the avatars")
	flag.Parse()

	osPort := os.Getenv("PORT")
	if osPort != "" {
		config.Port = osPort
	}

	return config
}

func main() {
	config := inflateDaemonConfig()

	if err := daemon.Run(config); err != nil {
		fmt.Printf("Error during daemon startup: %s\n", err)
	}
}
