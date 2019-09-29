package main

import (
	"back/daemon"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	var isLocalhost bool

	flag.BoolVar(&isLocalhost, "local", false, "Is it local mashine")
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

	if !isLocalhost {
		config.StaticDirPath = "/public"
		files, err := ioutil.ReadDir("/")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name(), file.IsDir())
		}
	}

	log.Println(http.Dir("public"))

	return config
}

func main() {
	config := inflateDaemonConfig()

	if err := daemon.Run(config); err != nil {
		fmt.Printf("Error during daemon startup: %s\n", err)
	}
}
