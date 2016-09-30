package mq

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/surgemq/message"
)

func readMessageBuffer(r io.Reader) (buf []byte, err error) {
	var l int
	var h [5]byte

	// fixed header
	// MQTT协议规定，最多允许４个字节表示剩余长度，每个字节最高位若为1，则表示还有后续字节存在
	for {
		n, err := r.Read(h[l : l+1])
		if err != nil {
			return nil, err
		}
		if n != 1 {
			return nil, errors.New("the number of bytes read is not 1")
		}
		l += 1

		if l > 1 && h[l-1] < 0x80 {
			break
		}

		if l > 5 {
			return nil, fmt.Errorf("4th byte of remaining length has continuation bit set")
		}
	}

	remlen, _ := binary.Uvarint(h[1:])
	buf = make([]byte, uint64(l)+remlen)
	copy(buf, h[:l])

	// variable header + payload
	for l < len(buf) {
		n, err := r.Read(buf[l:])
		if err != nil {
			return nil, err
		}
		l += n
	}

	return buf, nil
}

func readConnectMessage(r io.Reader) (*message.ConnectMessage, error) {
	buf, err := readMessageBuffer(r)
	if err != nil {
		return nil, err
	}

	msg := message.NewConnectMessage()
	if _, err := msg.Decode(buf); err != nil {
		return nil, err
	}
	return msg, nil
}

func decodeBuffer(buf []byte) (message.Message, error) {
	mtype := message.MessageType(buf[0] >> 4)
	msg, err := mtype.New()
	if err != nil {
		return nil, err
	}
	if _, err = msg.Decode(buf); err != nil {
		return nil, err
	}
	return msg, nil
}

func readMessage(r io.Reader) (message.Message, error) {
	buf, err := readMessageBuffer(r)
	if err != nil {
		return nil, err
	}
	return decodeBuffer(buf)
}

func writeMessage(w io.Writer, m message.Message) (err error) {
	buf := make([]byte, m.Len())
	if _, err = m.Encode(buf); err != nil {
		return err
	}
	if _, err = w.Write(buf); err != nil {
		return err
	}
	return nil
}
