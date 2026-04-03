package bcs

import (
	"encoding/binary"
	"math/big"
)

type Serializable interface {
	Serialize(s *Serializer)
}

type Serializer struct {
	buf []byte
}

func (s *Serializer) SerializeBool(v bool) {
	if v {
		s.buf = append(s.buf, 1)
	} else {
		s.buf = append(s.buf, 0)
	}
}

func (s *Serializer) SerializeU8(v uint8) {
	s.buf = append(s.buf, v)
}

func (s *Serializer) SerializeU16(v uint16) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	s.buf = append(s.buf, b...)
}

func (s *Serializer) SerializeU32(v uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	s.buf = append(s.buf, b...)
}

func (s *Serializer) SerializeU64(v uint64) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	s.buf = append(s.buf, b...)
}

func (s *Serializer) SerializeU128(v *big.Int) {
	b := make([]byte, 16)
	vb := v.Bytes()
	for i, byt := range vb {
		b[len(vb)-1-i] = byt
	}
	s.buf = append(s.buf, b...)
}

func (s *Serializer) SerializeU256(v *big.Int) {
	b := make([]byte, 32)
	vb := v.Bytes()
	for i, byt := range vb {
		b[len(vb)-1-i] = byt
	}
	s.buf = append(s.buf, b...)
}

func (s *Serializer) SerializeULEB128(v uint64) {
	for {
		b := byte(v & 0x7f)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		s.buf = append(s.buf, b)
		if v == 0 {
			break
		}
	}
}

func (s *Serializer) SerializeBytes(v []byte) {
	s.SerializeULEB128(uint64(len(v)))
	s.buf = append(s.buf, v...)
}

func (s *Serializer) SerializeFixedBytes(v []byte) {
	s.buf = append(s.buf, v...)
}

func (s *Serializer) SerializeString(v string) {
	s.SerializeBytes([]byte(v))
}

func (s *Serializer) ToBytes() []byte {
	out := make([]byte, len(s.buf))
	copy(out, s.buf)
	return out
}
