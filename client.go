package main

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Client struct {
	conn   *websocket.Conn
	source chan string
	quit   chan struct{}
}

func newClient(c *websocket.Conn) *Client {
	return &Client{
		conn: c,
		quit: globalQuit,
	}
}

// this reads from the web socket, but we don't do
// anything with that information here.
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		body := make(map[string]string)
		err := c.conn.ReadJSON(&body)
		if err != nil {
			log.Printf("read error: %s", err)
			clients.Remove(c)
			return
		} else {
			log.Printf("message body: %s", body)
		}
	}
}

// this writes messages from client.source to the websocket
// as those messages are recieved on the client's particular
// source channel
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	c.source = make(chan string, 10)
	for {
		select {
		case message := <-c.source:
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("write error: %s", err)
				clients.Remove(c)
				return
			}
		case <-c.quit:
			return // terminate the client
		}
	}
}

type Clients struct {
	sync.Mutex
	clients map[*Client]*Client
}

func newClients() *Clients {
	return &Clients{
		clients: make(map[*Client]*Client),
	}
}

func (cs *Clients) Add(c *Client) {
	cs.Lock()
	defer cs.Unlock()
	cs.clients[c] = c
}

func (cs *Clients) Remove(c *Client) {
	cs.Lock()
	defer cs.Unlock()
	delete(cs.clients, c)
}

func (cs *Clients) Iter(routine func(*Client)) {
	cs.Lock()
	defer cs.Unlock()
	log.Printf("sending to %d clients", len(cs.clients))
	for c := range cs.clients {
		routine(c)
	}
}
