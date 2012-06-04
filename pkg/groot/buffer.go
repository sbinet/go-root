package groot

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Buffer struct {
	order  binary.ByteOrder
	buffer []byte
	pos    int
	klen   uint32 // to compute refs (used in read_class, read_object)
}

func NewBuffer(buf []byte, order binary.ByteOrder, klen uint32) (b *Buffer, err error) {
	b = &Buffer{
		order:  order,
		buffer: buf[:],
		pos:    0,
		klen:   klen,
	}
	return b, err
}

func NewBufferFromKey(k Key) (b *Buffer, err error) {
	buf, err := k.Buffer()
	if err != nil {
		return
	}
	b = &Buffer{
		order: k.file.order,
		buffer: buf,
		pos: 0,
		klen: uint32(k.keysz),
	}
	return
}

func (b *Buffer) Buffer() []byte {
	return b.buffer[:]
}

func (b *Buffer) breader() breader {
	return breader{b.order}
}

func (b *Buffer) read_class_tag() (class string, err error) {
	br := b.breader()
	buf := bytes.NewBuffer(b.buffer[b.pos:])
	tag := br.ntou4(buf)
	println("tag:",tag)

	if tag == kNewClassTag {
		
	} else if (tag & kClassMask) != 0 {
		
	} else {
		println("**err** tag unknown case:", tag)
		err = fmt.Errorf("groot.Buffer: unknown tag case (%v)", tag)
		return
	}
	return
}

// EOF
