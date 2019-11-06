package handlers

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/log"
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
	hr "2019_2_Next_Level/internal/serverapi/server/Error/httpError"
	"2019_2_Next_Level/pkg/HttpTools"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type MailHandler struct {
	usecase mailbox.MailBoxUseCase
	resp    *HttpTools.Response
}

// NewMailHandler : sets handlers for specified routes (prefix = "/mail")
func NewMailHandler(router *mux.Router, usecase mailbox.MailBoxUseCase) {
	handler := MailHandler{usecase: usecase}
	handler.resp = (&HttpTools.Response{}).SetError(hr.DefaultResponse)

	router.HandleFunc("/send", handler.SendMail).Methods("POST")
	router.HandleFunc("/list", handler.GetMailList).Methods("GET")
	router.HandleFunc("/mail/{id}", handler.GetMail).Methods("GET")

}

func (h *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	email := model.Email{}
	err := HttpTools.StructFromBody(*r, &email)
	if err !=nil {
		log.Log().E(err)
		return
	}
	err = h.usecase.SendMail(email.From, email.To, email.Body)
	if err != nil {
		log.Log().E("Cannot send email")
	}
}

func (h *MailHandler) GetMailList(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := r.Header.Get("X-Login")
	type Request struct {
		Sort string
		Since int
		Count int
		Folder string
	}
	req := Request{"time-fresh-first", 1, 100, "incoming"}
	//err := HttpTools.StructFromBody(*r, &req)
	//if err != nil {
	//	resp.SetError(hr.BadParam)
	//	return
	//}


	list, err := h.usecase.GetMailList(login, req.Folder, req.Sort, req.Since, req.Count)
	if err != nil {
		resp.SetError(hr.BadParam)
		return
	}
	err = HttpTools.BodyFromStruct(w, list)
	if err != nil {
		resp.SetError(hr.BadParam)
		return
	}
}

func (h *MailHandler) GetMail(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := r.Header.Get("X-Login")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		resp.SetError(hr.BadParam)
		return
	}
	mail, err := h.usecase.GetMail(login, id)
	if err != nil {
		resp.SetError(hr.BadParam)
		return
	}
	data, _ := json.Marshal(mail)
	w.Write(data)
}