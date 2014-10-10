package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/mchmarny/ws-collector/backend"
	"github.com/mchmarny/ws-collector/common"
	"io"
	"log"
)

const (
	// WebSockets buffer
	channelBufSize = 100

	// Kafka type of backend
	backendTypeKAFKA = "kafka"

	// AMQP type of backend
	backendTypeAMQP = "amqp"

	// MQTT type of backend
	backendTypeMQTT = "mqtt"
)

var (
	maxId int = 0
	msg       = websocket.Message
)

type handler struct {
	id     int
	ws     *websocket.Conn
	server *broker
	ch     chan *interface{}
	doneCh chan bool
	sender chan *common.Message
}

type Backend interface {
	Config(uri string, args map[string]string)
	Start(in <-chan *common.Message)
}

func newClient(ws *websocket.Conn, s *broker) *handler {
	if ws == nil {
		panic("ws cannot be nil")
	}
	if s == nil {
		panic("server cannot be nil")
	}
	maxId++
	ch := make(chan *interface{}, channelBufSize)
	doneCh := make(chan bool)

	h := &handler{maxId, ws, s, ch, doneCh, make(chan *common.Message, 1)}

	var q Backend

	switch args.Backend.Type {
	case backendTypeKAFKA:
		q = backend.NewKafkaBackend()
	case backendTypeAMQP:
		q = backend.NewAMQPBackend()
	case backendTypeMQTT:
		q = backend.NewMQTTBackend()
	default:
		log.Panicf("invalid backend type: %s", args.Backend.Type)
	}

	log.Printf("configuring/starting %s backend...", args.Backend.Type)
	q.Config(args.Backend.URI, args.Backend.Args)
	go q.Start(h.sender)

	return h
}

func (c *handler) write(msg *interface{}) {
	select {
	case c.ch <- msg:
	default:
		c.server.del(c)
		err := fmt.Errorf("handler %d is disconnected.", c.id)
		c.server.err(err)
	}
}

func (c *handler) conn() *websocket.Conn { return c.ws }
func (c *handler) bone()                 { c.doneCh <- true }
func (c *handler) listen()               { c.listenRead() }
func (c *handler) listenRead() {
	log.Println("reading from handler")
	for {
		select {
		case <-c.doneCh:
			c.server.del(c)
			c.doneCh <- true
			return
		default:
			var m string
			err := msg.Receive(c.ws, &m)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.err(err)
			} else {
				//TODO: send it to queue here
				log.Print(m)
				c.sender <- common.NewMessage(m)
			}
		}
	}
}
