package http

import (
	hr "2019_2_Next_Level/internal/serverapi/server/Error/httpError"
	"2019_2_Next_Level/internal/support/log"
	"2019_2_Next_Level/internal/support/models"
	"2019_2_Next_Level/internal/support/usecase"
	"github.com/go-park-mail-ru/2019_2_Next_Level/pkg/HttpTools"
	"github.com/gorilla/mux"
	"net/http"
)

type SupportHandler struct {
	usecase usecase.Usecase
	newMessageChan map[string]chan interface{}
}

var NewMessageNotify *usecase.WebSocket

func NewSupportHandler(usecase_ usecase.Usecase) *SupportHandler {
	NewMessageNotify = usecase.NewWebSocket()
	return &SupportHandler{usecase: usecase_}
}

func (h *SupportHandler) InflateRouter(router *mux.Router) {
	h.newMessageChan = make(map[string]chan interface{}, 100)
	router.HandleFunc("/create", h.StartChat).Methods("POST")
	router.HandleFunc("/chat", h.GetChatList).Methods("GET")
	router.HandleFunc("/chat/{id}", h.GetChat).Methods("GET")
	router.HandleFunc("/chat/{id}/close", h.CloseChat).Methods("POST")
	router.HandleFunc("/chat/{id}/send", h.SendMessage).Methods("POST")
	router.HandleFunc("/chat/{id}/message/{message_id}", h.GetMessage).Methods("GET")
	router.HandleFunc("/newmessage", NewMessageNotify.ListenNewMessageNotify())
}

func (h *SupportHandler) StartChat(w http.ResponseWriter, r *http.Request) {
	resp := HttpTools.NewResponse(w)
	defer resp.Send()

	var req struct {
		Theme string `json:"theme"`
	}
	err := HttpTools.StructFromBody(*r, &req)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}

	login := h.getLogin(r)

	chatId, supporterName, err := h.usecase.StartChat(login, req.Theme)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}

	ans := struct {
		Status string `json:"status`
		ChatId interface{} `json:"chatId"`
		Supporter string `json:"supporterName"`
	}{
		Status:"OK",
		ChatId:chatId,
		Supporter:supporterName,
	}
	resp.SetContent(ans)
	return


}

func (h *SupportHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	resp := HttpTools.NewResponse(w)
	defer resp.Send()

	args := mux.Vars(r)
	chatId, ok := args["id"]
	if !ok {
		log.Log().E("No param id")
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	login := h.getLogin(r)
	chat, err := h.usecase.GetChat(login, chatId)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}

	ans := struct {
		Status string      `json:"status"`
		Chat   models.Chat `json:"chat"`
	}{
		Status:"OK", Chat:chat,
	}
	resp.SetContent(ans)
}

func (h *SupportHandler) GetChatList(w http.ResponseWriter, r *http.Request) {
	resp := HttpTools.NewResponse(w)
	defer resp.Send()
	login := h.getLogin(r)
	chats, err := h.usecase.GetChats(login)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}
	ans := struct {
		Status string       `json:"status"`
		Chats []models.Chat `json:"chats"`
	} {
		Status:"OK", Chats:chats,
	}
	resp.SetContent(ans)
}

func (h *SupportHandler) CloseChat(w http.ResponseWriter, r *http.Request) {
	resp := HttpTools.NewResponse(w)
	defer resp.Send()

	args := mux.Vars(r)
	chatId, ok := args["id"]
	if !ok {
		log.Log().E("No param id")
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	login := h.getLogin(r)

	err := h.usecase.CloseChat(login, chatId)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}

	resp.SetContent(hr.DefaultResponse)
}

func (h *SupportHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	resp := HttpTools.NewResponse(w)
	defer resp.Send()

	args := mux.Vars(r)
	chatId, ok := args["id"]
	if !ok {
		log.Log().E("No param id")
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	login := h.getLogin(r)

	var req models.Message
	err := HttpTools.StructFromBody(*r, &req)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}

	err = h.usecase.SendMessage(login, chatId, req)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}
	resp.SetContent(hr.DefaultResponse)
}

func (h *SupportHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	resp := HttpTools.NewResponse(w)
	defer resp.Send()

	args := mux.Vars(r)
	chatId, ok := args["id"]
	if !ok {
		log.Log().E("No param id")
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	messageId, ok := args["message_id"]
	if !ok {
		log.Log().E("No param message_id")
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	login := h.getLogin(r)

	mess, err := h.usecase.GetMessage(login, chatId, messageId)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}
	ans := struct{
		Status  string         `json:"status"`
		Message models.Message `json:"message"`
	}{
		Status: hr.OK, Message:mess,
	}
	resp.SetContent(ans)
}

func (h *SupportHandler) getLogin(r *http.Request) string {
	return r.Header.Get("X-Login")
}

