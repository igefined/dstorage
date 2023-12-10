package p2p

import "net"

type TCPPeer struct {
	conn           net.Conn
	isOutboundPeer bool
}

func NewTCPPeer(conn net.Conn, isOutbound bool) *TCPPeer {
	return &TCPPeer{
		conn:           conn,
		isOutboundPeer: isOutbound,
	}
}

func (p TCPPeer) Close() error {
	return p.conn.Close()
}
