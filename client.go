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
		log.Println("readPump closing")
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
		log.Println("writePump closing")
		c.conn.Close()
	}()

	for {
		select {
		case message := <-logchan:
			log.Printf("writing %s", message)
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("write error: %s", err)
				return
			}
		}
	}
}
