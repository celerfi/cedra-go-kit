package bcs

import (
	"math/big"
	"testing"
)

func TestBoolRoundTrip(t *testing.T) {
	for _, v := range []bool{true, false} {
		s := &Serializer{}
		s.SerializeBool(v)
		d := NewDeserializer(s.ToBytes())
		got := d.DeserializeBool()
		if got != v {
			t.Errorf("bool round-trip: want %v got %v", v, got)
		}
	}
}

func TestU8RoundTrip(t *testing.T) {
	cases := []uint8{0, 1, 127, 255}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeU8(v)
		d := NewDeserializer(s.ToBytes())
		if got := d.DeserializeU8(); got != v {
			t.Errorf("u8 round-trip: want %d got %d", v, got)
		}
	}
}

func TestU16RoundTrip(t *testing.T) {
	cases := []uint16{0, 1, 256, 65535}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeU16(v)
		d := NewDeserializer(s.ToBytes())
		if got := d.DeserializeU16(); got != v {
			t.Errorf("u16 round-trip: want %d got %d", v, got)
		}
	}
}

func TestU32RoundTrip(t *testing.T) {
	cases := []uint32{0, 1, 65536, 4294967295}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeU32(v)
		d := NewDeserializer(s.ToBytes())
		if got := d.DeserializeU32(); got != v {
			t.Errorf("u32 round-trip: want %d got %d", v, got)
		}
	}
}

func TestU64RoundTrip(t *testing.T) {
	cases := []uint64{0, 1, 1<<32 - 1, 1<<63 - 1, 1<<64 - 1}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeU64(v)
		d := NewDeserializer(s.ToBytes())
		if got := d.DeserializeU64(); got != v {
			t.Errorf("u64 round-trip: want %d got %d", v, got)
		}
	}
}

func TestU128RoundTrip(t *testing.T) {
	cases := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
		new(big.Int).Lsh(big.NewInt(1), 64),
		new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)),
	}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeU128(v)
		d := NewDeserializer(s.ToBytes())
		got := d.DeserializeU128()
		if got.Cmp(v) != 0 {
			t.Errorf("u128 round-trip: want %s got %s", v, got)
		}
	}
}

func TestU256RoundTrip(t *testing.T) {
	cases := []*big.Int{
		big.NewInt(0),
		big.NewInt(255),
		new(big.Int).Lsh(big.NewInt(1), 128),
		new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)),
	}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeU256(v)
		d := NewDeserializer(s.ToBytes())
		got := d.DeserializeU256()
		if got.Cmp(v) != 0 {
			t.Errorf("u256 round-trip: want %s got %s", v, got)
		}
	}
}

func TestULEB128RoundTrip(t *testing.T) {
	cases := []uint64{0, 1, 127, 128, 16383, 16384, 2097151, 268435455, 1<<63 - 1}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeULEB128(v)
		d := NewDeserializer(s.ToBytes())
		got := d.DeserializeULEB128()
		if got != v {
			t.Errorf("ULEB128 round-trip: want %d got %d", v, got)
		}
	}
}

func TestBytesRoundTrip(t *testing.T) {
	cases := [][]byte{
		{},
		{0x00},
		{0xff, 0xfe, 0xfd},
		make([]byte, 200),
	}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeBytes(v)
		d := NewDeserializer(s.ToBytes())
		got := d.DeserializeBytes()
		if len(got) != len(v) {
			t.Errorf("bytes round-trip length: want %d got %d", len(v), len(got))
			continue
		}
		for i := range v {
			if got[i] != v[i] {
				t.Errorf("bytes round-trip byte %d: want %d got %d", i, v[i], got[i])
			}
		}
	}
}

func TestStringRoundTrip(t *testing.T) {
	cases := []string{"", "hello", "cedra-go-kit", "0x1::cedra_coin::CedraCoin"}
	for _, v := range cases {
		s := &Serializer{}
		s.SerializeString(v)
		d := NewDeserializer(s.ToBytes())
		got := d.DeserializeString()
		if got != v {
			t.Errorf("string round-trip: want %q got %q", v, got)
		}
	}
}

func TestFixedBytes(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	s := &Serializer{}
	s.SerializeFixedBytes(data)
	d := NewDeserializer(s.ToBytes())
	got := d.DeserializeFixedBytes(5)
	for i, b := range data {
		if got[i] != b {
			t.Errorf("fixed bytes mismatch at %d: want %d got %d", i, b, got[i])
		}
	}
}

func TestDeserializerErrorOnShortRead(t *testing.T) {
	d := NewDeserializer([]byte{0x01})
	d.DeserializeU64()
	if d.Error() == nil {
		t.Error("expected error on short read, got nil")
	}
}

func TestMultipleValuesSequential(t *testing.T) {
	s := &Serializer{}
	s.SerializeU8(42)
	s.SerializeString("hello")
	s.SerializeU64(999)

	d := NewDeserializer(s.ToBytes())
	if v := d.DeserializeU8(); v != 42 {
		t.Errorf("want 42 got %d", v)
	}
	if v := d.DeserializeString(); v != "hello" {
		t.Errorf("want hello got %q", v)
	}
	if v := d.DeserializeU64(); v != 999 {
		t.Errorf("want 999 got %d", v)
	}
	if d.Error() != nil {
		t.Errorf("unexpected error: %v", d.Error())
	}
}
