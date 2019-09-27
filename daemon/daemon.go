package daemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
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
	// fmt.Fprintln(w, file)
	// w.Header().Set("Content-Type", "text/html")

	contentType := getFileType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(file)

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

	server := http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Millisecond,
		WriteTimeout: 10 * time.Millisecond,
	}

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
