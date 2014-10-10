package common

import (
	"encoding/json"
	"log"
	"time"
)

func NewMessage(msg string) *Message {
	return &Message{time.Now().UTC(), msg}
}

type Message struct {

	// when it was received
	On time.Time `json:"on"`

	// message boxy
	Body string `json:"body"`
}

func (m *Message) ToBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("unable to marshal: %v", err.Error())
	}
	return b
}
