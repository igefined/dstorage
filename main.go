package main

import (
	"fmt"
	"log"

	"github.com/igefined/dstorage/p2p"
)

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		Handshake:  p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
		OnPeer: func(peer p2p.Peer) error {
			return nil
		},
	}
	tr := p2p.NewTCPTransport(opts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Consume message: %+v\n", msg)
		}
	}()

	if err := tr.Listen(); err != nil {
		log.Fatal(err)
	}

	select {}
}
