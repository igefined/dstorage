package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: ":4300",
		Handshake:  func(peer Peer) error { return nil },
		Decoder:    DefaultDecoder{},
	}

	tr := NewTCPTransport(opts).(*TCPTransport)
	assert.Equal(t, opts.ListenAddr, tr.ListenAddr)

	err := tr.Listen()
	assert.NoError(t, err)
}
