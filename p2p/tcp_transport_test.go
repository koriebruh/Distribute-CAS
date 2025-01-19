package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTCPTransport(t *testing.T) {
	listenerAddr := ":4000"
	tr := NewTCPTransport(listenerAddr)
	assert.Equal(t, tr.listenAddress, listenerAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
