package server

import (
	authhandler "2019_2_Next_Level/internal/serverapi/server/Auth/http"
	authrepo "2019_2_Next_Level/internal/serverapi/server/Auth/repository"
	authusecase "2019_2_Next_Level/internal/serverapi/server/Auth/usecase"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/handlers"
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

	db "back/database"
)

func Run(externwg *sync.WaitGroup) error {
	// if conn == nil {
	// 	// create personal connction to db
	// }
	defer externwg.Done()
	fmt.Println("Starting daemon on port ", config.Conf.Port)

	// authUseCase := authusecase.NewAuthUsecase()

	db.Init()
	mainRouter := mux.NewRouter()
	router := mainRouter.PathPrefix("/api").Subrouter()

	InflateRouter(router)

	err := http.ListenAndServe(config.Conf.Port, router)
	return err
}

func InflateRouter(router *mux.Router) {
	router.Use(middleware.CorsMethodMiddleware()) // CORS for all requests
	router.Use(middleware.AccessLogMiddleware())

	authRouter := router.PathPrefix("/auth").Subrouter()
	// authRepo := authrepo.GetMock()
	authRepo, err := authrepo.GetPostgres()
	if err != nil {
		fmt.Println("Error during init Postgres", err)
		return
	}
	authUseCase := authusecase.NewAuthUsecase(&authRepo)
	authHandler := authhandler.NewAuthHandler(&authUseCase)
	authHandler.InflateRouter(authRouter)

	userRouter := router.PathPrefix("/profile").Subrouter()
	userRouter.Use(middleware.AuthentificationMiddleware(&authUseCase))
	userRepo, err := userrepo.GetPostgres()
	if err != nil {
		fmt.Println("Error during init Postgres", err)
		return
	}
	userUsecase := userusecase.NewUserUsecase(&userRepo)
	userHandler := userhandler.NewUserHandler(&userUsecase)
	userHandler.InflateRouter(userRouter)

	mailRouter := router.PathPrefix("/mail").Subrouter()
	mailRouter.Use(middleware.AuthentificationMiddleware(&authUseCase))
	handlers.NewMailHandler(mailRouter, &mailboxusecase.MailBoxUsecase{})

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
	})
}
