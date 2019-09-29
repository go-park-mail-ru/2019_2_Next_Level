package daemon

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	db "back/database"
)

type DataHandler struct {
}

func (h *DataHandler) getWorkDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
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
	outUser := (&UserOutput{}).FromUser(user)

	if len(config.AvatarDirPath) > 0 && config.AvatarDirPath[len(config.AvatarDirPath)-1] != '/' {
		config.AvatarDirPath = config.AvatarDirPath + "/"
	}
	outUser.AvaUrl = config.AvatarDirPath + db.GetAvaFilename(user)
	js, err := json.Marshal(outUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *DataHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest
	(&CorsHandler{}).preflightHandler(w, r)
	email, err := (&AuthHandler{}).CheckAuthorization(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(status)
		return
	}

	user, err := (&AuthHandler{}).parseUser(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		fmt.Println("Cannot get avatar", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(h.getWorkDirectory()+config.AvatarDirPath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Cannot create avatar file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	user.Email = email
	db.UpdateUser(user)
}

func (h *DataHandler) GetFront(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	fmt.Printf("Path: %s\n", r.URL.Path)
	path = config.StaticDirPath + path
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

func (h *DataHandler) GetPersonalFile(w http.ResponseWriter, r *http.Request) {
	log.Println("GetPersonalFile")
	status := http.StatusBadRequest
	(&CorsHandler{}).preflightHandler(w, r)
	_, err := (&AuthHandler{}).CheckAuthorization(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(status)
		return
	}
	h.GetFront(w, r)
}

func (h *DataHandler) GetOpenFile(w http.ResponseWriter, r *http.Request) {
	log.Println("get open file", r.URL.Path)
	path := r.URL.Path
	if res, _ := regexp.Match("^/open", []byte(path)); !res {
		path = "/open" + path
	}
	if res, _ := regexp.Match("/$", []byte(path)); res {
		path += "index.html"
	}
	r.URL.Path = path
	h.GetFront(w, r)
}

func (h *DataHandler) getFileType(filename string) string {
	type typePair struct {
		Type  string
		Value string
	}
	// textPrefix := "text/"
	types := []typePair{
		{"js", "text/javascript"},
		{"html", "text/html"},
		{"css", "text/css"},
		{"png", "image/png"},
	}

	for _, elem := range types {
		reg := fmt.Sprintf(`.%s$`, elem.Type)
		if res, _ := regexp.MatchString(reg, filename); res {
			return elem.Value
		}
	}
	return "text/plain"
}
