package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	conn *websocket.Conn
}

func newClient(c *websocket.Conn) *Client {
	return &Client{
		conn: c,
	}
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		body := make(map[string]string)
		err := c.conn.ReadJSON(&body)
		if err != nil {
			log.Printf("read error: %s", err)
			return
		}
		log.Printf("message body: %s", body)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message := <-logchan:
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("write error: %s", err)
				return
			}
		}
	}
}
