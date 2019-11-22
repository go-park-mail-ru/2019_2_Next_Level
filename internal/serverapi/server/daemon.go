package server

import (
	"2019_2_Next_Level/internal/Auth"
	"2019_2_Next_Level/internal/serverapi/log"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	authhandler "2019_2_Next_Level/internal/serverapi/server/Auth/http"
	authrepo "2019_2_Next_Level/internal/serverapi/server/Auth/repository"
	authusecase "2019_2_Next_Level/internal/serverapi/server/Auth/usecase"
	mailhandler "2019_2_Next_Level/internal/serverapi/server/MailBox/handlers"
	mailrepo "2019_2_Next_Level/internal/serverapi/server/MailBox/repository"
	mailboxusecase "2019_2_Next_Level/internal/serverapi/server/MailBox/usecase"
	userhandler "2019_2_Next_Level/internal/serverapi/server/User/http"
	userrepo "2019_2_Next_Level/internal/serverapi/server/User/repository"
	userusecase "2019_2_Next_Level/internal/serverapi/server/User/usecase"
	"2019_2_Next_Level/internal/serverapi/server/config"
	"2019_2_Next_Level/internal/serverapi/server/middleware"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

)

func Run(externwg *sync.WaitGroup) error {
	defer externwg.Done()
	log.Log().L("Starting daemon on port ", config.Conf.Port)

	mainRouter := mux.NewRouter()
	router := mainRouter.PathPrefix("/api").Subrouter()

	InflateRouter(router)

	//mainRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(config.Conf.StaticDir))))
	//staticHandler := http.StripPrefix("/", http.FileServer(http.Dir(config.Conf.StaticDir)))
	//mainRouter.PathPrefix("/").Handler(middleware.StaticMiddleware()(staticHandler))

	staticHandler := http.Handler(http.FileServer(http.Dir(config.Conf.StaticDir)))
	//mainRouter.PathPrefix("/").Handler(http.FileServer(http.Dir(config.Conf.StaticDir)))
	mainRouter.PathPrefix("/").Handler(middleware.StaticMiddleware()(staticHandler))
	//mainRouter.PathPrefix("/").
	//	Subrouter().
	//	Handle("/", middleware.StaticMiddleware()(http.FileServer(http.Dir(config.Conf.StaticDir))))

	err := http.ListenAndServe(config.Conf.Port, mainRouter)
	return err
}

func InflateRouter(router *mux.Router) {
	router.Use(middleware.CorsMethodMiddleware()) // CORS for all requests
	router.Use(middleware.AccessLogMiddleware())
	router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	authRouter := router.PathPrefix("/auth").Subrouter()
	authUseCase := InitHttpAuth(authRouter)

	userRouter := router.PathPrefix("/profile").Subrouter()
	userRouter.Use(middleware.AuthentificationMiddleware(authUseCase))
	InitHttpUser(userRouter)

	mailRouter := router.PathPrefix("/messages").Subrouter()
	mailRouter.Use(middleware.AuthentificationMiddleware(authUseCase))
	InitHTTPMail(mailRouter)

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
	})

}

func InitHttpAuth(router *mux.Router) auth.Usecase {
	authRepo, err := authrepo.GetPostgres()
	if err != nil {
		log.Log().E("Error during init Postgres", err)
		return nil
	}
	authClient := Auth.AuthClient{}
	err =authClient.Init("0.0.0.0", ":6000")
	if err != nil {
		log.Log().E("Cannot init auth: ", err)
	}
	authUseCase := authusecase.NewAuthUsecase(&authRepo, &authClient)
	authHandler := authhandler.NewAuthHandler(&authUseCase)
	authHandler.InflateRouter(router)
	return &authUseCase
}

func InitHttpUser(router *mux.Router) {
	userRepo, err := userrepo.GetPostgres()
	if err != nil {
		log.Log().E("Error during init Postgres", err)
		return
	}
	authClient := Auth.AuthClient{}
	err =authClient.Init("0.0.0.0", ":6000")
	if err != nil {
		log.Log().E("Cannot init auth: ", err)
	}
	userUsecase := userusecase.NewUserUsecase(&userRepo, &authClient)
	userHandler := userhandler.NewUserHandler(&userUsecase)
	userHandler.InflateRouter(router)
}

func InitHTTPMail(router *mux.Router) {
	repo, err := mailrepo.GetPostgres()
	if err != nil {
		log.Log().E("Error during init Postgres", err)
		return
	}
	mailUsecase := mailboxusecase.NewMailBoxUsecase(&repo)
	handler := mailhandler.NewMailHandler(mailUsecase)
	handler.InflateRouter(router)
	//handlers.NewMailHandler(router, &mailUsecase)
}
