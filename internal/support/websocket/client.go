package websocket

import (
	"2019_2_Next_Level/internal/support/models"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
)

type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *models.Message
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	return &Client{ws: ws, server: server}
}
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}
func (c *Client) listenWrite() {
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			//c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *Client) listenRead() {
	for {
		select {

		// receive done request
		case <-c.doneCh:
			//c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg models.Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				//c.server.Err(err)
				fmt.Println("Error")
			} else {
				fmt.Println("SendAll")
				//c.server.SendAll(&msg)
			}
		}
	}
}