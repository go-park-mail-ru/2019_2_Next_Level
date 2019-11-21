package main

import (
	"2019_2_Next_Level/internal/Auth"
	"2019_2_Next_Level/internal/Auth/log"
	"2019_2_Next_Level/internal/Auth/repository"
	pb "2019_2_Next_Level/internal/Auth/service"
	"2019_2_Next_Level/pkg/wormhole"
	"google.golang.org/grpc"
)

const (
	port = ":6000"
)

func main() {
	log.Log().SetPrefix("Auth")
	repo := repository.NewPostgresRepo()
	err := repo.Init("postgres", "postgres", "0.0.0.0", "5432", "nextlevel")
	if err != nil {
		log.Log().E("Cannot init repo: ", err)
		return
	}
	authWorker := Auth.NewAuthServer(repo)
	hole := wormhole.Wormhole{}
	err = hole.RunServer(port, func(server *grpc.Server) {
		pb.RegisterAuthServer(server, authWorker)
		log.Log().L("Auth started on port: ", port)
	})
	if err != nil {
		log.Log().E("Error during work of Auth service: ", err)
	}

}
