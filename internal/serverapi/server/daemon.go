package server

import (
	"2019_2_Next_Level/internal/model"
	authusecase "2019_2_Next_Level/internal/serverapi/server/Auth/usecase"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/handlers"
	mailboxusecase "2019_2_Next_Level/internal/serverapi/server/MailBox/usecase"
	"2019_2_Next_Level/internal/serverapi/server/auth/repo"
	"2019_2_Next_Level/internal/serverapi/server/middleware"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	db "back/database"
)

func Run(externwg *sync.WaitGroup, conn *model.Connection) error {
	if conn == nil {
		// create personal connction to db
	}
	defer externwg.Done()
	fmt.Println("Starting daemon on port ", Conf.Port)
	repo.

	authUserCase := authusecase.AuthUsecase.SetRepo(repo.NewPostgres(conn))
	mailUseCase := mail.MailUseCase{repository.Postgres{conn}}

	db.Init()

	// authMux := mux.NewRouter()
	router := mux.NewRouter()

	private := router.PathPrefix("/mail").Subrouter()
	// userMux := mux.NewRouter()
	handlers.NewMailHandler(private, &mailboxusecase.MailBoxUsecase{})
	router.Use(mux.CORSMethodMiddleware(router))

	private.Use(middleware.AuthentificationMiddleware(&authusecase.AuthUsecase{}))
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
	})

	err := http.ListenAndServe(Conf.Port, router)
	return err
}
