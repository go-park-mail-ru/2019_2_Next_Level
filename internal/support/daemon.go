package support

import (
	"2019_2_Next_Level/internal/Auth"
	http2 "2019_2_Next_Level/internal/support/http"
	repository "2019_2_Next_Level/internal/support/repository"
	"2019_2_Next_Level/internal/support/usecase"

	"2019_2_Next_Level/internal/support/log"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

const (
	Port = ":7000"
)
const (
	dbuser = "postgres"
	dbpass = "postgres"
	dbhost = "0.0.0.0"
	dbport = "5432"
	dbname = "nextlevel"
	port = "5000"
)

func Run(externwg *sync.WaitGroup) error {
	defer externwg.Done()
	log.Log().L("Starting daemon on port ", Port)

	mainRouter := mux.NewRouter()
	router := mainRouter.PathPrefix("/api").Subrouter()
	InflateRouter(router)

	err := http.ListenAndServe(Port, mainRouter)
	log.Log().E("Error of http.ListenAndServe(): ", err)
	return err
}

func InflateRouter(router *mux.Router) {
	//router.Use(middleware.CorsMethodMiddleware()) // CORS for all requests
	//router.Use(middleware.AccessLogMiddleware())
	router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	supportRouter := router.PathPrefix("/support").Subrouter()
	InitHttpSupport(supportRouter)


}

func InitHttpSupport(router *mux.Router) usecase.Usecase {
	repo := repository.NewPostgresRepository(nil)
	err := repo.Init(dbuser, dbpass, dbhost, dbport, dbname)
	if err != nil {
		log.Log().E("Cannot init repo: ", err)
	}
	authClient := Auth.AuthClient{}
	err =authClient.Init("0.0.0.0", ":6000")
	if err != nil {
		log.Log().E("Cannot init auth: ", err)
	}
	UseCase := usecase.NewSupportUsecase(repo)
	authHandler := http2.NewSupportHandler(UseCase)
	authHandler.InflateRouter(router)
	return UseCase
}

