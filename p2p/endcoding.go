package p2p

import (
	"bufio"
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(reader io.Reader, msg *RPC) error
}

type GOBDecoder struct{}

func (g GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct {
}

func (d DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	reader := bufio.NewReader(r)         // Gunakan bufio untuk membaca per baris
	line, err := reader.ReadString('\n') // Baca sampai newline (\n)
	if err != nil {
		return err
	}

	msg.Payload = []byte(line)

	return nil
}
