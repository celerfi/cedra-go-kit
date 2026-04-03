package bcs

import (
	"encoding/binary"
	"errors"
	"math/big"
)

type Deserializer struct {
	data   []byte
	offset int
	err    error
}

func NewDeserializer(data []byte) *Deserializer {
	return &Deserializer{data: data}
}

func (d *Deserializer) Error() error {
	return d.err
}

func (d *Deserializer) Remaining() int {
	return len(d.data) - d.offset
}

func (d *Deserializer) read(n int) ([]byte, error) {
	if d.err != nil {
		return nil, d.err
	}
	if d.offset+n > len(d.data) {
		d.err = errors.New("bcs: unexpected end of data")
		return nil, d.err
	}
	b := d.data[d.offset : d.offset+n]
	d.offset += n
	return b, nil
}

func (d *Deserializer) DeserializeBool() bool {
	b, err := d.read(1)
	if err != nil {
		return false
	}
	return b[0] != 0
}

func (d *Deserializer) DeserializeU8() uint8 {
	b, err := d.read(1)
	if err != nil {
		return 0
	}
	return b[0]
}

func (d *Deserializer) DeserializeU16() uint16 {
	b, err := d.read(2)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint16(b)
}

func (d *Deserializer) DeserializeU32() uint32 {
	b, err := d.read(4)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint32(b)
}

func (d *Deserializer) DeserializeU64() uint64 {
	b, err := d.read(8)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(b)
}

func (d *Deserializer) DeserializeU128() *big.Int {
	b, err := d.read(16)
	if err != nil {
		return new(big.Int)
	}
	rev := make([]byte, 16)
	for i, byt := range b {
		rev[15-i] = byt
	}
	return new(big.Int).SetBytes(rev)
}

func (d *Deserializer) DeserializeU256() *big.Int {
	b, err := d.read(32)
	if err != nil {
		return new(big.Int)
	}
	rev := make([]byte, 32)
	for i, byt := range b {
		rev[31-i] = byt
	}
	return new(big.Int).SetBytes(rev)
}

func (d *Deserializer) DeserializeULEB128() uint64 {
	var result uint64
	var shift uint
	for {
		b, err := d.read(1)
		if err != nil {
			return 0
		}
		result |= uint64(b[0]&0x7f) << shift
		if b[0]&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 64 {
			d.err = errors.New("bcs: ULEB128 overflow")
			return 0
		}
	}
	return result
}

func (d *Deserializer) DeserializeBytes() []byte {
	length := d.DeserializeULEB128()
	if d.err != nil {
		return nil
	}
	b, err := d.read(int(length))
	if err != nil {
		return nil
	}
	out := make([]byte, length)
	copy(out, b)
	return out
}

func (d *Deserializer) DeserializeFixedBytes(n int) []byte {
	b, err := d.read(n)
	if err != nil {
		return nil
	}
	out := make([]byte, n)
	copy(out, b)
	return out
}

func (d *Deserializer) DeserializeString() string {
	return string(d.DeserializeBytes())
}
