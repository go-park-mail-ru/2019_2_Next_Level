package usecase

import (
	"fmt"
	"net/http"
	gorillaWebsocket "github.com/gorilla/websocket"
)

type WebSocket struct {
	upgrader gorillaWebsocket.Upgrader
	chans map[string]chan interface{}
}

func NewWebSocket() *WebSocket {
	s := WebSocket{upgrader: gorillaWebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}}
	s.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	s.chans = make(map[string]chan interface{})
	return &s
}
func (s *WebSocket) ListenNewMessageNotify(input chan interface{}) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		conn, _ := s.upgrader.Upgrade(w, r, nil)
		// Принять клиента
		var token []byte
		var err error
		for{
			_, token, err = conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			if _, ok := s.chans[string(token)]; !ok {
				s.chans[string(token)] = make(chan interface{}, 100)
			}
		}
		select{
		case res:=<-s.chans[string(token)]:
			if err := conn.WriteMessage(gorillaWebsocket.TextMessage, []byte(res.(string))); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (s *WebSocket) NewNotify(userToken string, id string) {
	s.chans[userToken]<-id
}