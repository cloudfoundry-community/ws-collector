package common

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func NewAuth(token string) *Auth {
	a := &Auth{token: token}
	a.enabled = len(token) > 0
	return a
}

type Auth struct {
	enabled bool
	token   string
}

func (a *Auth) Valid(req *http.Request) bool {
	if !a.enabled {
		return true
	}
	auths, _ := req.Header["Authorization"]
	if len(auths) != 1 {
		log.Println("missing Authorization")
		return false
	}
	tokens := strings.Split(auths[0], " ")
	if len(tokens) != 2 || tokens[0] != "Bearer" {
		log.Printf("invalid auth type: %s", tokens)
		return false
	}
	raw, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		log.Printf("unable to decode token: %s", tokens[1])
		return false
	}
	token := string(raw)
	log.Printf("validating token: %s", token)
	return a.token == token
}
