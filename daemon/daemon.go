package daemon

import (
	"fmt"
	"io/ioutil"
	"log"
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
	files, erro := ioutil.ReadDir(".")
	if erro != nil {
		log.Fatal(erro)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	config = *cfg
	fmt.Println("Starting daemon on port ", cfg.Port)
	db.Init()

	router := mux.NewRouter()
	authApi := &AuthHandler{}
	dataApi := &DataHandler{}
	corsApi := &CorsHandler{}
	router.HandleFunc("/api/auth/signup", authApi.Register).Methods("POST")
	router.HandleFunc("/api/auth/signin", authApi.Login).Methods("POST")
	router.HandleFunc("/api/auth/signout", authApi.Logout).Methods("POST")
	router.HandleFunc("/api/profile/get", dataApi.GetProfile).Methods("GET")
	router.HandleFunc("/api/profile/{field}", dataApi.UpdateProfile).Methods("POST")

	// router.HandleFunc("/signup", authApi.Register).Methods("POST")
	// router.HandleFunc("/signin", authApi.Login).Methods("POST")
	// router.HandleFunc("/signout", authApi.Logout).Methods("POST")
	// router.HandleFunc("/profile", dataApi.GetProfile).Methods("GET")
	// router.HandleFunc("/profile/{field}", dataApi.UpdateProfile).Methods("POST")

	router.PathPrefix("/private").HandlerFunc(dataApi.GetPersonalFile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(dataApi.GetOpenFile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(corsApi.preflightHandler).Methods("OPTIONS")

	err := http.ListenAndServe(":"+cfg.Port, router)
	return err
}
