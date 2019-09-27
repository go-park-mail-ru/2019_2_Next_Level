package daemon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	db "../database"
)

type Config struct {
	Port         int
	BackendPath  string
	FrontendPath string
}
type Myerror struct {
}

var config Config

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	// fmt.Printf("Path: %s\n", r.URL.Path)
	path = config.FrontendPath + path
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Cannot open file %s: %s\n", path, err)
		return
	}
	fmt.Printf("Give %s file\n", path)

	contentType := getFileType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(file)

}

func getProfileHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest
	// w.WriteHeader(status)
	// defer func() {
	// 	w.WriteHeader(status)
	// 	fmt.Printf("Status: %d\n", status)
	// }()
	session, err := r.Cookie("user-id")
	if err != nil {
		fmt.Println("Unauthorized user")
		w.WriteHeader(status)
		return
	}
	fmt.Printf("Cookie: %s\n", session.Value)
	email, err := db.GetUserEmailBySession(session.Value)
	if err != nil {
		fmt.Println("Unauthorized user")
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
	status = http.StatusAccepted
	// w.WriteHeader(status)
}

func getFileType(filename string) string {
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

func Run(cfg *Config) error {
	config = *cfg
	fmt.Printf("Starting daemon on port %d\n", cfg.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/profile", getProfileHandler)

	server := http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Millisecond,
		WriteTimeout: 10 * time.Millisecond,
	}
	db.Init()
	server.ListenAndServe()
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Cannot start listening port %d", cfg.Port)
		return err
	}

	// http.HandleFunc("/", handler)

	// port := ":" + strconv.Itoa(cfg.Port)

	// if err := http.ListenAndServe(port, nil); err != nil {
	// 	log.Printf("Cannot start listening port %d", cfg.Port)
	// 	return err
	// }

	return nil
}
