package daemon

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	db "../database"
)

type Config struct {
	Port          int
	BackendPath   string
	FrontendPath  string
	FrontendUrl   string
	AvatarDirPath string
}

var config Config

func Run(cfg *Config) error {
	config = *cfg
	fmt.Printf("Starting daemon on port %d\n", cfg.Port)
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
	router.PathPrefix("/").HandlerFunc(dataApi.GetPersonalFile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(corsApi.preflightHandler).Methods("OPTIONS")

	err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router)
	return err
}
