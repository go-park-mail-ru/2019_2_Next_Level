package daemon

import (
	"encoding/json"
	"fmt"
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
	email, err := (&AuthHandler{}).CheckAuthorization(r)
	if err != nil {
		fmt.Println(err)
		(&Error{ErrorNoPermission}).Send(&w)
		return
	}
	user, err := db.GetUserByEmail(email)
	if err != nil {
		(&Error{ErrorNoPermission}).Send(&w)
		return
	}
	fmt.Println(user)
	outUser := (&UserOutput{}).FromUser(user)

	// if len(config.AvatarDirPath) > 0 && config.AvatarDirPath[len(config.AvatarDirPath)-1] != '/' {
	// 	config.AvatarDirPath = config.AvatarDirPath + "/"
	// }
	// outUser.AvaUrl = config.AvatarDirPath + db.GetAvaFilename(user)
	js, err := json.Marshal(outUser)
	if err != nil {
		(&Error{ErrorInternal}).Send(&w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *DataHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	(&CorsHandler{}).preflightHandler(w, r)
	email, err := (&AuthHandler{}).CheckAuthorization(r)
	if err != nil {
		fmt.Println(err)
		(&Error{ErrorNoPermission}).Send(&w)
		return
	}
	userInput := UserInput{}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &userInput)

	user, _ := db.GetUserByEmail(email)
	if userInput.Name != "" {
		user.Name = userInput.Name
	}
	if userInput.Password != "" {
		user.Password = userInput.Password
	}
	fmt.Println(user.Name)

	// switch vars["field"] {
	// case "change_name":
	// 	user.Name = r.PostForm["name"][0]
	// 	log.Println("Change name to ", user.Name)
	// 	break
	// default:
	// 	log.Println("Wrong param: ", vars["field"])
	// }
	db.UpdateUser(user)

	// user, err := (&AuthHandler{}).parseUser(r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	// file, handler, err := r.FormFile("avatar")
	// if err != nil {
	// 	fmt.Println("Cannot get avatar", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// defer file.Close()

	// f, err := os.OpenFile(h.getWorkDirectory()+config.AvatarDirPath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	fmt.Println("Cannot create avatar file", err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer f.Close()
	// io.Copy(f, file)

	// user.Email = email
	// db.UpdateUser(user)
}

func (h *DataHandler) GetFront(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path += configuration.FileForFolder
	}
	fmt.Printf("Path: %s\n", r.URL.Path)
	path = configuration.StaticDir + path
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Cannot open file %s: %s\n", path, err)
		(&Error{"Cannot open file " + path}).Send(&w)
		return
	}
	fmt.Printf("Give %s file\n", path)

	contentType := h.getFileType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(file)
}

func (h *DataHandler) GetPersonalFile(w http.ResponseWriter, r *http.Request) {
	log.Println("GetPersonalFile")
	(&CorsHandler{}).preflightHandler(w, r)
	_, err := (&AuthHandler{}).CheckAuthorization(r)
	if err != nil {
		fmt.Println(err)
		(&Error{ErrorNoPermission}).Send(&w)
		return
	}
	h.GetFront(w, r)
}

func (h *DataHandler) GetOpenFile(w http.ResponseWriter, r *http.Request) {
	log.Println("get open file", r.URL.Path)
	path := r.URL.Path
	if res, _ := regexp.Match("^/open", []byte(path)); !res {
		path = "/" + configuration.OpenDir + path
	}
	if res, _ := regexp.Match("/$", []byte(path)); res {
		path += configuration.FileForFolder
	}
	r.URL.Path = path
	h.GetFront(w, r)
}

func (h *DataHandler) getFileType(filename string) string {
	type typePair struct {
		Type  string
		Value string
	}
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
