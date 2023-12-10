package p2p

import (
	"errors"
	"fmt"
	"net"
)

type (
	TCPTransportOpts struct {
		ListenAddr string
		Handshake  HandshakeFunc
		Decoder    Decoder
		OnPeer     func(Peer) error
	}

	TCPTransport struct {
		listener net.Listener
		TCPTransportOpts
		rpcCh chan *RPC
	}
)

func NewTCPTransport(opts TCPTransportOpts) Transport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan *RPC),
	}
}

func (t *TCPTransport) Listen() (err error) {
	t.listener, err = net.Listen("tcp", t.ListenAddr)

	go t.startAccept()

	return
}

func (t *TCPTransport) Consume() <-chan *RPC {
	return t.rpcCh
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
	var err error

	defer func() {
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err = t.Handshake(peer); err != nil {
		fmt.Printf("TCP handshake error: %s", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			fmt.Printf("TCP OnPeer error: %s", err)
			return
		}
	}

	msg := &RPC{}
	for {
		err = t.Decoder.Decode(conn, msg)
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP can't read: %s", err)
			continue

		}

		msg.From = conn.RemoteAddr()

		t.rpcCh <- msg
	}
}
