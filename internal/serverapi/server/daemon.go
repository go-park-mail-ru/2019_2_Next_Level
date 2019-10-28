package server

import (
	authhandler "2019_2_Next_Level/internal/serverapi/server/Auth/http"
	authrepo "2019_2_Next_Level/internal/serverapi/server/Auth/repository"
	authusecase "2019_2_Next_Level/internal/serverapi/server/Auth/usecase"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/handlers"
	mailboxusecase "2019_2_Next_Level/internal/serverapi/server/MailBox/usecase"
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
	fmt.Println("Starting daemon on port ", Conf.Port)

	// authUseCase := authusecase.NewAuthUsecase()

	db.Init()

	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	InflateRouter(router)

	err := http.ListenAndServe(Conf.Port, router)
	return err
}

func InflateRouter(router *mux.Router) {
	router.Use(middleware.AccessLogMiddleware())
	router.Use(mux.CORSMethodMiddleware(router)) // CORS for all requests

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRepo := authrepo.GetMock()
	authUseCase := authusecase.NewAuthUsecase(&authRepo)
	authHandler := authhandler.NewAuthHandler(&authUseCase)
	authHandler.InflateRouter(authRouter)

	mailRouter := router.PathPrefix("/mail").Subrouter()
	mailRouter.Use(middleware.AuthentificationMiddleware(&authUseCase))
	handlers.NewMailHandler(mailRouter, &mailboxusecase.MailBoxUsecase{})

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
	})
}
