package unpack

import "encoding/binary"

type ReadBuff []byte

func (b *ReadBuff) Uint8() uint8 {
	v := (*b)[0]
	*b = (*b)[1:]
	return v
}

func (b *ReadBuff) Uint16() uint16 {
	v := binary.LittleEndian.Uint16(*b)
	*b = (*b)[2:]
	return v
}

func (b *ReadBuff) Uint32() uint32 {
	v := binary.LittleEndian.Uint32(*b)
	*b = (*b)[4:]
	return v
}

func (b *ReadBuff) Uint64() uint64 {
	v := binary.LittleEndian.Uint64(*b)
	*b = (*b)[8:]
	return v
}

func (b *ReadBuff) Sub(n int) ReadBuff {
	b2 := (*b)[:n]
	*b = (*b)[n:]
	return b2
}
