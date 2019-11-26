package main

import (
	pb "2019_2_Next_Level/generated/Auth/service"
	"2019_2_Next_Level/internal/Auth"
	localconf "2019_2_Next_Level/internal/Auth/config"
	"2019_2_Next_Level/internal/Auth/log"
	"2019_2_Next_Level/internal/Auth/repository"
	"2019_2_Next_Level/pkg/config"
	"2019_2_Next_Level/pkg/wormhole"
	"flag"
	"google.golang.org/grpc"
)


func main() {
	log.Log().SetPrefix("Auth")

	err := initializeConfig()
	if err != nil {
		log.Log().E(err)
		return
	}

	repo := repository.NewPostgresRepo()
	err = repo.Init(
		localconf.Conf.DB.User, localconf.Conf.DB.Password, localconf.Conf.DB.Host, localconf.Conf.DB.Port, localconf.Conf.DB.DBName)
	if err != nil {
		log.Log().E("Cannot init repo: ", err)
		return
	}
	authWorker := Auth.NewAuthServer(repo)
	hole := wormhole.Wormhole{}
	port := localconf.Conf.AuthPort
	err = hole.RunServer(port, func(server *grpc.Server) {
		pb.RegisterAuthServer(server, authWorker)
		log.Log().L("Auth started on port: ", port)
	})
	if err != nil {
		log.Log().E("Error during work of Auth service: ", err)
	}

}

func initializeConfig() error {
	configFilename := flag.String("config", "default", "Path to config file")
	dbUser := flag.String("dbuser", "", "User for database")
	dbPassword := flag.String("dbpass", "", "Password for database")
	flag.Parse()

	return config.Inflate(*configFilename, &localconf.Conf, *dbUser, *dbPassword)
}

