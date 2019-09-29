package daemon

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	// db "../database"
	db "back/database"
)

type Config struct {
	Port          string
	BackendPath   string
	FrontendPath  string
	FrontendUrl   string
	AvatarDirPath string
	StaticDirPath string
}

var config Config

func Run(cfg *Config) error {
	config = *cfg
	fmt.Println("Starting daemon on port ", cfg.Port)
	db.Init()

	router := mux.NewRouter()
	authApi := &AuthHandler{}
	dataApi := &DataHandler{}
	corsApi := &CorsHandler{}
	// router.HandleFunc("/auth.signup", authApi.Register).Methods("POST")
	// router.HandleFunc("/auth.signin", authApi.Login).Methods("POST")
	// router.HandleFunc("/settings/profile", dataApi.GetProfile).Methods("GET")
	// router.HandleFunc("/settings/profile", dataApi.UpdateProfile).Methods("POST")
	router.HandleFunc("/signup", authApi.Register).Methods("POST")
	router.HandleFunc("/signin", authApi.Login).Methods("POST")
	router.HandleFunc("/profile", dataApi.GetProfile).Methods("GET")
	router.HandleFunc("/profile", dataApi.UpdateProfile).Methods("POST")
	// router.PathPrefix("/").HandlerFunc(dataApi.GetFront).Methods("GET")
	router.PathPrefix("/private").HandlerFunc(dataApi.GetPersonalFile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(dataApi.GetOpenFile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(corsApi.preflightHandler).Methods("OPTIONS")

	err := http.ListenAndServe(":"+cfg.Port, router)
	return err
}
