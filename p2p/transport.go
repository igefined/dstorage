package p2p

type Peer interface {
	Close() error
}

type Transport interface {
	Listen() error
	Consume() <-chan *RPC
}
