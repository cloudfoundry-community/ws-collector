package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/mchmarny/ws-collector/common"
	"log"
	"net/http"
	"os"
)

type broker struct {
	clients map[int]*handler
	addCh   chan *handler
	delCh   chan *handler
	doneCh  chan bool
	errCh   chan error
	authVal *common.Auth
}

func new() *broker {
	clients := make(map[int]*handler, 5)
	addCh := make(chan *handler, 5)
	delCh := make(chan *handler)
	doneCh := make(chan bool)
	errCh := make(chan error)
	authV := common.NewAuth(args.Server.Token)

	return &broker{
		clients,
		addCh,
		delCh,
		doneCh,
		errCh,
		authV,
	}
}

func (s *broker) add(c *handler) {
	if s.authVal.Valid(c.ws.Request()) {
		s.addCh <- c
	}
}
func (s *broker) del(c *handler) { s.delCh <- c }
func (s *broker) done()          { s.doneCh <- true }
func (s *broker) err(err error)  { s.errCh <- err }
func (s *broker) listen() {
	log.Println("broker listening...")
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
		handler := newClient(ws, s)
		s.add(handler)
		handler.listen()
	}

	onRequest := func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(onConnected)}
		s.ServeHTTP(w, req)
	}

	http.HandleFunc(args.Server.Root, onRequest)

	log.Println("handler created")

	for {
		select {
		case c := <-s.addCh:
			log.Println("new handler added")
			s.clients[c.id] = c
			log.Printf("connected clients: %d", len(s.clients))
		case c := <-s.delCh:
			log.Println("handler deleted")
			delete(s.clients, c.id)
		case err := <-s.errCh:
			log.Println("error:", err.Error())
		case <-s.doneCh:
			return
		}
	}
}

func showHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Requested: %s", r.URL.Path)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	b := new()
	go b.listen()
	a := fmt.Sprintf("%s:%d", args.Server.Host, args.Server.Port)
	http.HandleFunc("/", showHome)
	log.Printf("config: %v", args)
	log.Fatal(http.ListenAndServe(a, nil))
}
