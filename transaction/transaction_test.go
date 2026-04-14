package transaction

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/bcs"
)

func makeFixedFeePayerRawTxn(t *testing.T) (*RawTransaction, *FeePayerRawTransaction, *account.Ed25519Account, *account.Ed25519Account) {
	t.Helper()

	sender, err := account.NewEd25519AccountFromHex("0x1111111111111111111111111111111111111111111111111111111111111111")
	if err != nil {
		t.Fatalf("sender account: %v", err)
	}
	feePayer, err := account.NewEd25519AccountFromHex("0x2222222222222222222222222222222222222222222222222222222222222222")
	if err != nil {
		t.Fatalf("fee payer account: %v", err)
	}
	recipient, _ := account.AccountAddressFromHex("0xdeadbeef")
	moduleAddr, _ := account.AccountAddressFromHex("0x1")

	rawTxn := &RawTransaction{
		Sender:         sender.Address(),
		SequenceNumber: 7,
		Payload: &EntryFunction{
			Module:   ModuleID{Address: moduleAddr, Name: "cedra_account"},
			Function: "transfer",
			TypeArgs: []TypeTag{},
			Args: [][]byte{
				SerializeAddressArg(recipient),
				SerializeU64Arg(4242),
			},
		},
		MaxGasAmount:            200_000,
		GasUnitPrice:            100,
		ExpirationTimestampSecs: 1_700_000_000,
		ChainID:                 2,
	}

	return rawTxn, rawTxn.WithFeePayer(account.AccountAddress{}), sender, feePayer
}

func TestSerializeU64Arg(t *testing.T) {
	b := SerializeU64Arg(1000)
	if len(b) != 8 {
		t.Errorf("u64 arg should be 8 bytes, got %d", len(b))
	}
	d := bcs.NewDeserializer(b)
	if v := d.DeserializeU64(); v != 1000 {
		t.Errorf("want 1000 got %d", v)
	}
}

func TestSerializeU8Arg(t *testing.T) {
	b := SerializeU8Arg(42)
	if len(b) != 1 {
		t.Errorf("u8 arg should be 1 byte, got %d", len(b))
	}
	d := bcs.NewDeserializer(b)
	if v := d.DeserializeU8(); v != 42 {
		t.Errorf("want 42 got %d", v)
	}
}

func TestSerializeBoolArg(t *testing.T) {
	bt := SerializeBoolArg(true)
	bf := SerializeBoolArg(false)
	if bt[0] != 1 {
		t.Errorf("true should serialize to 1, got %d", bt[0])
	}
	if bf[0] != 0 {
		t.Errorf("false should serialize to 0, got %d", bf[0])
	}
}

func TestSerializeAddressArg(t *testing.T) {
	addr, _ := account.AccountAddressFromHex("0x1")
	b := SerializeAddressArg(addr)
	if len(b) != 32 {
		t.Errorf("address arg should be 32 bytes, got %d", len(b))
	}
	if b[31] != 1 {
		t.Errorf("last byte should be 1, got %d", b[31])
	}
}

func TestSerializeBytesArg(t *testing.T) {
	data := []byte{0xde, 0xad, 0xbe, 0xef}
	b := SerializeBytesArg(data)
	d := bcs.NewDeserializer(b)
	got := d.DeserializeBytes()
	if hex.EncodeToString(got) != hex.EncodeToString(data) {
		t.Errorf("bytes arg round-trip mismatch")
	}
}

func TestSerializeStringArg(t *testing.T) {
	s := "cedra"
	b := SerializeStringArg(s)
	d := bcs.NewDeserializer(b)
	got := d.DeserializeString()
	if got != s {
		t.Errorf("string arg round-trip: want %q got %q", s, got)
	}
}

func TestSerializeU128Arg(t *testing.T) {
	v := new(big.Int).Lsh(big.NewInt(1), 64)
	b := SerializeU128Arg(v)
	if len(b) != 16 {
		t.Errorf("u128 arg should be 16 bytes, got %d", len(b))
	}
	d := bcs.NewDeserializer(b)
	got := d.DeserializeU128()
	if got.Cmp(v) != 0 {
		t.Errorf("u128 arg round-trip: want %s got %s", v, got)
	}
}

func TestEntryFunctionSerialize(t *testing.T) {
	addr, _ := account.AccountAddressFromHex("0x1")
	ef := &EntryFunction{
		Module:   ModuleID{Address: addr, Name: "cedra_account"},
		Function: "transfer",
		TypeArgs: []TypeTag{},
		Args: [][]byte{
			SerializeAddressArg(addr),
			SerializeU64Arg(1_000_000),
		},
	}
	s := &bcs.Serializer{}
	ef.Serialize(s)
	b := s.ToBytes()
	if len(b) == 0 {
		t.Error("serialized entry function should not be empty")
	}
}

func TestRawTransactionSigningMessageLength(t *testing.T) {
	sender, _ := account.AccountAddressFromHex("0x1")
	moduleAddr, _ := account.AccountAddressFromHex("0x1")
	rawTxn := &RawTransaction{
		Sender:         sender,
		SequenceNumber: 0,
		Payload: &EntryFunction{
			Module:   ModuleID{Address: moduleAddr, Name: "cedra_account"},
			Function: "transfer",
			TypeArgs: []TypeTag{},
			Args: [][]byte{
				SerializeAddressArg(sender),
				SerializeU64Arg(1000),
			},
		},
		MaxGasAmount:            200_000,
		GasUnitPrice:            100,
		ExpirationTimestampSecs: 9999999999,
		ChainID:                 2,
	}

	s := &bcs.Serializer{}
	rawTxn.Serialize(s)
	txnBytes := s.ToBytes()
	if len(txnBytes) == 0 {
		t.Error("serialized raw transaction should not be empty")
	}
}

func TestRawTransactionSigningMessageDeterministic(t *testing.T) {
	sender, _ := account.AccountAddressFromHex("0xdeadbeef")
	moduleAddr, _ := account.AccountAddressFromHex("0x1")
	build := func() []byte {
		rawTxn := &RawTransaction{
			Sender:         sender,
			SequenceNumber: 5,
			Payload: &EntryFunction{
				Module:   ModuleID{Address: moduleAddr, Name: "coin"},
				Function: "transfer",
				TypeArgs: []TypeTag{TypeTagStruct{
					Address: moduleAddr,
					Module:  "cedra_coin",
					Name:    "CedraCoin",
				}},
				Args: [][]byte{
					SerializeAddressArg(sender),
					SerializeU64Arg(500),
				},
			},
			MaxGasAmount:            100_000,
			GasUnitPrice:            150,
			ExpirationTimestampSecs: 1700000000,
			ChainID:                 1,
		}
		s := &bcs.Serializer{}
		rawTxn.Serialize(s)
		return s.ToBytes()
	}

	b1 := build()
	b2 := build()
	if hex.EncodeToString(b1) != hex.EncodeToString(b2) {
		t.Error("raw transaction serialization is not deterministic")
	}
}

func TestFeePayerRawTransactionSerialize(t *testing.T) {
	_, feePayerRawTxn, _, _ := makeFixedFeePayerRawTxn(t)

	s := &bcs.Serializer{}
	feePayerRawTxn.Serialize(s)

	const wantHex = "01147e4d3a5b10eaed2a93536e284c23096dfcea9ac61f0a8420e5d01fbd8f0ea807000000000000000200000000000000000000000000000000000000000000000000000000000000010d63656472615f6163636f756e74087472616e7366657200022000000000000000000000000000000000000000000000000000000000deadbeef089210000000000000400d030000000000640000000000000000f1536500000000020700000000000000000000000000000000000000000000000000000000000000010a63656472615f636f696e094365647261436f696e00000000000000000000000000000000000000000000000000000000000000000000"
	if got := hex.EncodeToString(s.ToBytes()); got != wantHex {
		t.Fatalf("fee-payer raw txn mismatch\nwant: %s\ngot:  %s", wantHex, got)
	}
}

func TestTypeTagSerialization(t *testing.T) {
	tags := []TypeTag{
		TypeTagBool{},
		TypeTagU8{},
		TypeTagU16{},
		TypeTagU32{},
		TypeTagU64{},
		TypeTagU128{},
		TypeTagU256{},
		TypeTagAddress{},
		TypeTagSigner{},
		TypeTagVector{Element: TypeTagU8{}},
	}
	for _, tag := range tags {
		s := &bcs.Serializer{}
		tag.serializeTypeTag(s)
		if len(s.ToBytes()) == 0 {
			t.Errorf("type tag %T serialized to empty bytes", tag)
		}
	}
}

func TestParseTypeTag(t *testing.T) {
	cases := []struct {
		input    string
		wantType string
	}{
		{"bool", "*transaction.TypeTagBool"},
		{"u8", "*transaction.TypeTagU8"},
		{"u64", "*transaction.TypeTagU64"},
		{"address", "*transaction.TypeTagAddress"},
		{"vector<u8>", "*transaction.TypeTagVector"},
		{"0x1::cedra_coin::CedraCoin", "*transaction.TypeTagStruct"},
	}
	for _, c := range cases {
		tag := ParseTypeTag(c.input)
		if tag == nil {
			t.Errorf("ParseTypeTag(%q) returned nil", c.input)
		}
	}
}
