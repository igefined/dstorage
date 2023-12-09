package p2p

import (
	"fmt"
	"net"
	"sync"
)

type (
	TCPPeer struct {
		conn           net.Conn
		isOutboundPeer bool
	}

	TCPTransportOpts struct {
		ListenAddr string
		Handshake  HandshakeFunc
		Decoder    Decoder
	}

	TCPTransport struct {
		listener net.Listener
		TCPTransportOpts

		mu    sync.RWMutex
		peers map[net.Addr]Peer
	}
)

func NewTCPPeer(conn net.Conn, isOutbound bool) *TCPPeer {
	return &TCPPeer{
		conn:           conn,
		isOutboundPeer: isOutbound,
	}
}

func NewTCPTransport(opts TCPTransportOpts) Transport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) Listen() (err error) {
	t.listener, err = net.Listen("tcp", t.ListenAddr)

	go t.startAccept()

	return
}

func (t *TCPTransport) startAccept() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %+v\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.Handshake(peer); err != nil {
		fmt.Printf("TCP handshake error: %s", err)
		conn.Close()
		return
	}

	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP can't read: %s", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("TCP incoming msg: %+v\n", msg)
	}
}
