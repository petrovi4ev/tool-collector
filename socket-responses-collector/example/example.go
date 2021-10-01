package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	collector "github.com/BitMedia-IO/tool-collector/socket-responses-collector"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

// todo fix Bad request response
func main() {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	log.Println(s.URL)

	ctx, cancel := context.WithCancel(context.Background())
	req := &http.Request{
		URL: &url.URL{Host: s.URL},
	}
	wsMessageCollector := collector.New(s, req, echo)
	wsMessageCollector.Run(ctx)
	time.Sleep(5 * time.Second)
	cancel()

	log.Printf("Collected messages: %+v", wsMessageCollector.Messages())
}
