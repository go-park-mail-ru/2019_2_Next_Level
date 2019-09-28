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
type Myerror struct {
}

var config Config

type CorsHandler struct {
}

func (h *CorsHandler) preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Access-Control-Allow-Origin", config.FrontendUrl)
	headers.Add("Access-Control-Allow-Credentials", "true")
	headers.Add("Access-Control-Allow-Headers", "Content-Type")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

type UserInput struct {
	Name     string
	Email    string
	Password string
}
type UserOutput struct {
	Name       string
	Sirname    string
	MiddleName string
	Email      string
	AvaUrl     string
}

func (u *UserOutput) FromUser(dbuser db.User) UserOutput {
	user := UserOutput{
		Name:       dbuser.Name,
		Sirname:    dbuser.Sirname,
		MiddleName: dbuser.MiddleName,
		Email:      dbuser.Email,
	}
	return user
}

func (input *UserInput) ToUser() db.User {
	return db.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
}

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

	http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router)
	return nil
}
