package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4300"
	tr := NewTCPTransport(listenAddr).(*TCPTransport)
	assert.Equal(t, listenAddr, tr.listenAddress)

	err := tr.Listen()
	assert.NoError(t, err)
}
