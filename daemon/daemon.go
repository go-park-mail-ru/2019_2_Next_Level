package daemon

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"back/config"
	db "back/database"
)

var configuration config.Config

func Run(cfg *config.Config) error {
	configuration = *cfg
	fmt.Println("Starting daemon on port ", cfg.Port)
	db.Init()

	router := mux.NewRouter()
	authApi := &AuthHandler{}
	dataApi := &DataHandler{}
	corsApi := &CorsHandler{}
	router.HandleFunc("/api/auth/signup", authApi.Register).Methods("POST")
	router.HandleFunc("/api/auth/signin", authApi.Login).Methods("POST")
	router.HandleFunc("/api/auth/signout", authApi.Logout).Methods("GET")
	router.HandleFunc("/api/profile/get", dataApi.GetProfile).Methods("GET")
	router.HandleFunc("/api/profile/edit", dataApi.UpdateProfile).Methods("POST")
	router.PathPrefix("/private").HandlerFunc(dataApi.GetPersonalFile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(dataApi.GetOpenFile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(corsApi.preflightHandler).Methods("OPTIONS")

	err := http.ListenAndServe(":"+cfg.Port, router)
	return err
}
