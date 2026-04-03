package account

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/celerfi/cedra-go-kit/bcs"
)

type Account interface {
	Address() AccountAddress
	AuthKey() []byte
	PublicKeyBytes() []byte
	SignTransaction(signingMessage []byte) ([]byte, error)
}

type AccountAddress [32]byte

func NewAccountAddress(b []byte) (AccountAddress, error) {
	var addr AccountAddress
	if len(b) > 32 {
		return addr, errors.New("account: address too long")
	}
	copy(addr[32-len(b):], b)
	return addr, nil
}

func AccountAddressFromHex(s string) (AccountAddress, error) {
	s = strings.TrimPrefix(s, "0x")
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return AccountAddress{}, err
	}
	return NewAccountAddress(b)
}

func (a AccountAddress) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a AccountAddress) Bytes() []byte {
	b := make([]byte, 32)
	copy(b, a[:])
	return b
}

func (a AccountAddress) Serialize(s *bcs.Serializer) {
	s.SerializeFixedBytes(a[:])
}

func (a AccountAddress) String() string {
	return a.Hex()
}
