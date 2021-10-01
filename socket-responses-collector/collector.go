package socket_responses_collector

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

type Collector interface {
	Run(ctx context.Context)
	Messages() ResponsesMap
	Clean()
}

type ResponsesCollector struct {
	server   *httptest.Server
	request  *http.Request
	messages ResponseMessages
	handler  http.HandlerFunc
	ws       *websocket.Conn
}

func New(server *httptest.Server, request *http.Request, handler http.HandlerFunc) Collector {
	messages := ResponseMessages{M: make(ResponsesMap), mx: &sync.Mutex{}}

	return &ResponsesCollector{
		server:   server,
		request:  request,
		messages: messages,
		handler:  handler,
	}
}

func (collector *ResponsesCollector) Run(_ context.Context) {
	err := collector.dial()
	if err != nil {
		log.Fatalf("ws dial error: %s", err.Error())
	}

	go func() {
		for {
			_, p, err := collector.ws.ReadMessage()
			if err != nil {
				log.Fatalf("ws ReadMessage error: %s", err.Error())
			}
			collector.messages.Load(string(p))
		}
	}()

	collector.sendRequest()
}

func (collector *ResponsesCollector) Clean() {
	collector.messages.mx.Lock()
	defer collector.messages.mx.Unlock()

	collector.messages.M = make(ResponsesMap, 0)
}

func (collector *ResponsesCollector) dial() error {
	u := "ws" + strings.TrimPrefix(collector.server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	collector.ws = ws

	return err
}

func (collector *ResponsesCollector) sendRequest() {
	rr := httptest.NewRecorder()
	collector.handler.ServeHTTP(rr, collector.request)
	collector.ws.WriteMessage(websocket.TextMessage, []byte("jujuju"))
	collector.messages.Store(collector.request.URL.RawQuery, rr.Body.String())
}
