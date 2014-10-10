package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {

	assert.NotNil(t, args, "nil config args")
	assert.NotNil(t, args.Backend, "nil backend")
	assert.NotNil(t, args.Backend.Type, "nil type")
	assert.NotNil(t, args.Backend.URI, "nil uri")
	assert.NotNil(t, args.Backend.Args, "nil args")
	assert.NotNil(t, args.Server, "nil server")
	assert.NotNil(t, args.Server.Host, "nil host")
	assert.NotEqual(t, args.Server.Port, 0, "nil Port")
	assert.NotNil(t, args.Server.Root, "nil root")
	assert.NotNil(t, args.Server.Token, "nil root")

}
