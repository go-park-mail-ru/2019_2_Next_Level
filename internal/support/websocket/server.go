package websocket

import (
	"2019_2_Next_Level/internal/support/models"
	"golang.org/x/net/websocket"
	"net/http"
)

type Server struct {
	pattern   string
	messages  []*models.Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *models.Message
	doneCh    chan bool
	errCh     chan error
}

func (s *Server) AddClient(client *Client) {
	s.clients[10] = client
}

func NewServer(pattern string) *Server {
	return &Server{pattern: pattern}
}

func (s *Server) Add(c *Client) {
	s.clients[100] = c
}

func (s *Server) Listen() {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
}

//func (s *Server) Listen() {
//	for {
//		select {
//
//		// Add new a client
//		case c := <-s.addCh:
//			s.clients[c.id] = c
//			s.sendPastMessages(c)
//
//		// del a client
//		case c := <-s.delCh:
//			delete(s.clients, c.id)
//
//		// broadcast message for all clients
//		case msg := <-s.sendAllCh:
//			s.messages = append(s.messages, msg)
//			s.sendAll(msg)
//
//		case err := <-s.errCh:
//			log.Println("Error:", err.Error())
//
//		case <-s.doneCh:
//			return
//		}
//	}
//}