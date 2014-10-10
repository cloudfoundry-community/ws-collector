package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntParsers(t *testing.T) {

	i := ParseInt("1", 0)
	assert.Equal(t, i, 1)

	i = ParseInt("", 3)
	assert.Equal(t, i, 3)

}
