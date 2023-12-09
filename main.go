package main

import (
	"log"

	"github.com/igefined/dstorage/p2p"
)

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		Handshake:  p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(opts)

	if err := tr.Listen(); err != nil {
		log.Fatal(err)
	}

	select {}
}
