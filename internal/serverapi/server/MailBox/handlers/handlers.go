package handlers

import (
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
	"net/http"

	"github.com/gorilla/mux"
)

type MailHandler struct {
	mailusecase mailbox.MailBoxUseCase
}

// NewMailHandler : sets handlers for specified routes (prefix = "/mail")
func NewMailHandler(router *mux.Router, usecase mailbox.MailBoxUseCase) {
	handler := MailHandler{mailusecase: usecase}

	router.HandleFunc("/send", handler.SendMail).Methods("POST")
}

func (h *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	h.mailusecase.SendMail("aa@mail.ru", "ivan@yandex.ru", "Hello")
}
