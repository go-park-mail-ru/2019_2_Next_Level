package main

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/serverapi"
	"fmt"
)

func main() {
	fmt.Println("API Server started. Hello!")
	var a post.Sender
	a = &serverapi.QueueClient{}
	serverapi.SetQueue(a)
	serverapi.Run()

	// config.Configuration = config.Config{}
	// Configurate(&config.Configuration)

	// if err := daemon.Run(&config.Configuration); err != nil {
	// 	fmt.Printf("Error during daemon startup: %s\n", err)
	// }
}

// func Configurate(conf *config.Config) error {
// 	var isLocalhost bool

// 	flag.BoolVar(&isLocalhost, "local", false, "Is it local mashine")
// 	flag.Parse()
// 	if err := (*conf).Inflate(); err != nil {
// 		log.Println("Cannot read config: ", err)
// 		return errors.New("Cannot read config")
// 	}
// 	if !isLocalhost {
// 		osPort := os.Getenv("PORT")
// 		if osPort != "" {
// 			(*conf).Port = osPort
// 		}
// 	}
// 	if isLocalhost {
// 		(*conf).SelfURL = "http://localhost:" + (*conf).Port
// 	}
// 	return nil

// }
