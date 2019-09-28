package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	uuid "github.com/google/uuid"

	db "../database"
)

type Config struct {
	Port         int
	BackendPath  string
	FrontendPath string
	FrontendUrl  string
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

type AuthHandler struct {
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	handler := &CorsHandler{}
	handler.preflightHandler(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Cannot get body")
		http.Error(w, err.Error(), 500)
		return
	}
	user := UserInput{}
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println("Error during parse profile", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println("No such a user")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dbUser.Password != user.Password {
		fmt.Println("Wrong pasword")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Right user")
	a.Authorize(&w, &dbUser)
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	handler := &CorsHandler{}
	handler.preflightHandler(w, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Cannot get body")
		http.Error(w, err.Error(), 500)
		return
	}

	userInput := UserInput{}
	if err := json.Unmarshal(body, &userInput); err != nil {
		fmt.Println("Error during parse profile", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("User")
	user := userInput.ToUser()
	db.SetUser(user)

	a.Authorize(&w, &user)
}

func (a *AuthHandler) CheckAuthorization(r *http.Request) (string, error) {
	session, err := r.Cookie("user-id")
	if err != nil {
		return "", errors.New("No cookie")
	}
	fmt.Printf("Cookie: %s\n", session.Value)
	email, err := db.GetUserEmailBySession(session.Value)
	if err != nil {
		return "", errors.New("Wrong session key")
	}
	return email, nil
}

func (a *AuthHandler) Authorize(w *http.ResponseWriter, user *db.User) {
	out, _ := uuid.NewUUID()
	db.RegisterNewSession(out.String(), user.Email)
	cookie := http.Cookie{
		Name:    "user-id",
		Value:   out.String(),
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(*w, &cookie)
}

type DataHandler struct {
}

func (h *DataHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	(&CorsHandler{}).preflightHandler(w, r)

	status := http.StatusBadRequest

	email, err := (&AuthHandler{}).CheckAuthorization(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(status)
		return
	}
	user, err := db.GetUserByEmail(email)
	if err != nil {
		fmt.Printf("Cannot get user: %s\n", err)
		w.WriteHeader(status)
		return
	}
	fmt.Println(user)
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *DataHandler) GetFront(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	fmt.Printf("Path: %s\n", r.URL.Path)
	path = config.FrontendPath + path
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Cannot open file %s: %s\n", path, err)
		return
	}
	fmt.Printf("Give %s file\n", path)

	contentType := h.getFileType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(file)
}

func (h *DataHandler) getFileType(filename string) string {
	type typePair struct {
		Type  string
		Value string
	}
	textPrefix := "text/"
	types := []typePair{
		{"js", "javascript"},
		{"html", "html"},
		{"css", "css"},
	}

	for _, elem := range types {
		reg := fmt.Sprintf(`.%s$`, elem.Type)
		if res, _ := regexp.MatchString(reg, filename); res {
			return textPrefix + elem.Value
		}
	}
	return textPrefix + "plain"
}

type UserInput struct {
	Name     string
	Email    string
	Password string
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
	router.HandleFunc("/signup", authApi.Register).Methods("POST")
	router.HandleFunc("/signin", authApi.Login).Methods("POST")
	router.HandleFunc("/profile", dataApi.GetProfile).Methods("GET")
	router.PathPrefix("/").HandlerFunc(dataApi.GetFront).Methods("GET")
	router.PathPrefix("/").HandlerFunc(corsApi.preflightHandler).Methods("OPTIONS")

	http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router)
	return nil
}
