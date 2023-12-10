package p2p

import "io"

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type DefaultDecoder struct{}

func (d DefaultDecoder) Decode(reader io.Reader, rpc *RPC) error {
	buf := make([]byte, 1024)

	n, err := reader.Read(buf)
	if err != nil {
		return err
	}

	rpc.Payload = buf[:n]

	return nil
}
