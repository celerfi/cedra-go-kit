package transaction

import (
	"math/big"
	"strings"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/bcs"
)

const (
	typeTagBool    = 0
	typeTagU8      = 1
	typeTagU64     = 2
	typeTagU128    = 3
	typeTagAddress = 4
	typeTagSigner  = 5
	typeTagVector  = 6
	typeTagStruct  = 7
	typeTagU16     = 8
	typeTagU32     = 9
	typeTagU256    = 10
)

type TypeTag interface {
	serializeTypeTag(s *bcs.Serializer)
}

type TypeTagBool struct{}
type TypeTagU8 struct{}
type TypeTagU16 struct{}
type TypeTagU32 struct{}
type TypeTagU64 struct{}
type TypeTagU128 struct{}
type TypeTagU256 struct{}
type TypeTagAddress struct{}
type TypeTagSigner struct{}
type TypeTagVector struct{ Element TypeTag }
type TypeTagStruct struct {
	Address    account.AccountAddress
	Module     string
	Name       string
	TypeParams []TypeTag
}

func (TypeTagBool) serializeTypeTag(s *bcs.Serializer)    { s.SerializeULEB128(typeTagBool) }
func (TypeTagU8) serializeTypeTag(s *bcs.Serializer)      { s.SerializeULEB128(typeTagU8) }
func (TypeTagU16) serializeTypeTag(s *bcs.Serializer)     { s.SerializeULEB128(typeTagU16) }
func (TypeTagU32) serializeTypeTag(s *bcs.Serializer)     { s.SerializeULEB128(typeTagU32) }
func (TypeTagU64) serializeTypeTag(s *bcs.Serializer)     { s.SerializeULEB128(typeTagU64) }
func (TypeTagU128) serializeTypeTag(s *bcs.Serializer)    { s.SerializeULEB128(typeTagU128) }
func (TypeTagU256) serializeTypeTag(s *bcs.Serializer)    { s.SerializeULEB128(typeTagU256) }
func (TypeTagAddress) serializeTypeTag(s *bcs.Serializer) { s.SerializeULEB128(typeTagAddress) }
func (TypeTagSigner) serializeTypeTag(s *bcs.Serializer)  { s.SerializeULEB128(typeTagSigner) }

func (t TypeTagVector) serializeTypeTag(s *bcs.Serializer) {
	s.SerializeULEB128(typeTagVector)
	t.Element.serializeTypeTag(s)
}

func (t TypeTagStruct) serializeTypeTag(s *bcs.Serializer) {
	s.SerializeULEB128(typeTagStruct)
	t.Address.Serialize(s)
	s.SerializeString(t.Module)
	s.SerializeString(t.Name)
	s.SerializeULEB128(uint64(len(t.TypeParams)))
	for _, tp := range t.TypeParams {
		tp.serializeTypeTag(s)
	}
}

func ParseTypeTag(s string) TypeTag {
	switch s {
	case "bool":
		return TypeTagBool{}
	case "u8":
		return TypeTagU8{}
	case "u16":
		return TypeTagU16{}
	case "u32":
		return TypeTagU32{}
	case "u64":
		return TypeTagU64{}
	case "u128":
		return TypeTagU128{}
	case "u256":
		return TypeTagU256{}
	case "address":
		return TypeTagAddress{}
	case "signer":
		return TypeTagSigner{}
	}
	if strings.HasPrefix(s, "vector<") && strings.HasSuffix(s, ">") {
		inner := s[7 : len(s)-1]
		return TypeTagVector{Element: ParseTypeTag(inner)}
	}
	parts := strings.SplitN(s, "::", 3)
	if len(parts) == 3 {
		addr, _ := account.AccountAddressFromHex(parts[0])
		name := parts[2]
		typeParams := []TypeTag{}
		if idx := strings.Index(name, "<"); idx != -1 {
			inner := name[idx+1 : len(name)-1]
			name = name[:idx]
			for _, p := range strings.Split(inner, ",") {
				typeParams = append(typeParams, ParseTypeTag(strings.TrimSpace(p)))
			}
		}
		return TypeTagStruct{Address: addr, Module: parts[1], Name: name, TypeParams: typeParams}
	}
	return TypeTagU64{}
}

type ModuleID struct {
	Address account.AccountAddress
	Name    string
}

type EntryFunction struct {
	Module   ModuleID
	Function string
	TypeArgs []TypeTag
	Args     [][]byte
}

func (e *EntryFunction) Serialize(s *bcs.Serializer) {
	e.Module.Address.Serialize(s)
	s.SerializeString(e.Module.Name)
	s.SerializeString(e.Function)
	s.SerializeULEB128(uint64(len(e.TypeArgs)))
	for _, ta := range e.TypeArgs {
		ta.serializeTypeTag(s)
	}
	s.SerializeULEB128(uint64(len(e.Args)))
	for _, arg := range e.Args {
		s.SerializeBytes(arg)
	}
}

type RawTransaction struct {
	Sender                  account.AccountAddress
	SequenceNumber          uint64
	Payload                 *EntryFunction
	MaxGasAmount            uint64
	GasUnitPrice            uint64
	ExpirationTimestampSecs uint64
	ChainID                 uint8
	FeePayerCurrency        TypeTag
}

func defaultFeePayerCurrency() TypeTag {
	addr, _ := account.AccountAddressFromHex("0x1")
	return TypeTagStruct{Address: addr, Module: "cedra_coin", Name: "CedraCoin"}
}

func (r *RawTransaction) Serialize(s *bcs.Serializer) {
	r.Sender.Serialize(s)
	s.SerializeU64(r.SequenceNumber)
	s.SerializeULEB128(2)
	r.Payload.Serialize(s)
	s.SerializeU64(r.MaxGasAmount)
	s.SerializeU64(r.GasUnitPrice)
	s.SerializeU64(r.ExpirationTimestampSecs)
	s.SerializeU8(r.ChainID)
	currency := r.FeePayerCurrency
	if currency == nil {
		currency = defaultFeePayerCurrency()
	}
	currency.serializeTypeTag(s)
}

func SerializeBoolArg(v bool) []byte {
	s := &bcs.Serializer{}
	s.SerializeBool(v)
	return s.ToBytes()
}

func SerializeU8Arg(v uint8) []byte {
	s := &bcs.Serializer{}
	s.SerializeU8(v)
	return s.ToBytes()
}

func SerializeU64Arg(v uint64) []byte {
	s := &bcs.Serializer{}
	s.SerializeU64(v)
	return s.ToBytes()
}

func SerializeU128Arg(v *big.Int) []byte {
	s := &bcs.Serializer{}
	s.SerializeU128(v)
	return s.ToBytes()
}

func SerializeAddressArg(a account.AccountAddress) []byte {
	s := &bcs.Serializer{}
	a.Serialize(s)
	return s.ToBytes()
}

func SerializeBytesArg(v []byte) []byte {
	s := &bcs.Serializer{}
	s.SerializeBytes(v)
	return s.ToBytes()
}

func SerializeStringArg(v string) []byte {
	s := &bcs.Serializer{}
	s.SerializeString(v)
	return s.ToBytes()
}
