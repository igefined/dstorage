package p2p

import "io"

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type DefaultDecoder struct{}

func (d DefaultDecoder) Decode(reader io.Reader, msg *Message) error {
	buf := make([]byte, 1024)

	n, err := reader.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]

	return nil
}
