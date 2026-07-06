package main

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
)

type Decoder struct {
	r      io.Reader
	buf    []byte
	char   byte
	cursor int
}

var decoderPool = sync.Pool{
	New: func() interface{} {
		return &Decoder{}
	},
}

func NewDecoder(r io.Reader) *Decoder {
	dec := decoderPool.Get().(*Decoder)
	dec.Reset(r)
	return dec
}

func (d *Decoder) Reset(r io.Reader) {
	d.r = r
	d.buf = nil
	d.char = 0
	d.cursor = 0
}

func (d *Decoder) Release() {
	d.Reset(nil)
	decoderPool.Put(d)
}

func (d *Decoder) Decode(v interface{}) error {
	if d.buf == nil {
		var err error
		d.buf, err = io.ReadAll(d.r)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(d.buf, v)
}

func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewReader(data))
	defer dec.Release()
	return dec.Decode(v)
}
