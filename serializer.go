package chamux

import (
	"bytes"
	"encoding/gob"
)

type Serializer interface {
	Serialize(f *Frame) ([]byte, error)
	Deserialize(f []byte) (*Frame, error)
}

type Gob struct{}

func (g Gob) Serialize(f *Frame) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(f)
	return buf.Bytes(), err
}

func (g Gob) Deserialize(f []byte) (*Frame, error) {
	frame := &Frame{}
	var buf bytes.Buffer
	if _, err := buf.Write(f); err != nil {
		return frame, err
	}
	err := gob.NewDecoder(&buf).Decode(frame)
	return frame, err
}
