package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/celerfi/cedra-go-kit/bcs"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"golang.org/x/crypto/sha3"
)

type Secp256k1PrivateKey struct {
	key *secp256k1.PrivateKey
}

type Secp256k1PublicKey struct {
	key *secp256k1.PublicKey
}

type Secp256k1Signature struct {
	bytes [64]byte
}

func GenerateSecp256k1PrivateKey() (*Secp256k1PrivateKey, error) {
	key, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	return &Secp256k1PrivateKey{key: key}, nil
}

func NewSecp256k1PrivateKeyFromBytes(b []byte) (*Secp256k1PrivateKey, error) {
	if len(b) != 32 {
		return nil, errors.New("secp256k1: private key must be 32 bytes")
	}
	key := secp256k1.PrivKeyFromBytes(b)
	return &Secp256k1PrivateKey{key: key}, nil
}

func NewSecp256k1PrivateKeyFromHex(s string) (*Secp256k1PrivateKey, error) {
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return NewSecp256k1PrivateKeyFromBytes(b)
}

func (k *Secp256k1PrivateKey) PublicKey() *Secp256k1PublicKey {
	return &Secp256k1PublicKey{key: k.key.PubKey()}
}

func (k *Secp256k1PrivateKey) Sign(message []byte) (*Secp256k1Signature, error) {
	h := sha3.New256()
	h.Write(message)
	digest := h.Sum(nil)

	sig := ecdsa.SignCompact(k.key, digest, false)
	// SignCompact returns [recovery_flag(1)] + [R(32)] + [S(32)]
	if len(sig) != 65 {
		return nil, errors.New("secp256k1: unexpected signature length")
	}

	var out Secp256k1Signature
	copy(out.bytes[:32], sig[1:33])  // R
	copy(out.bytes[32:], sig[33:65]) // S

	// enforce low-S
	var s secp256k1.ModNScalar
	s.SetByteSlice(out.bytes[32:])
	if s.IsOverHalfOrder() {
		s.Negate()
		s.PutBytesUnchecked(out.bytes[32:])
	}

	return &out, nil
}

func (k *Secp256k1PrivateKey) Bytes() []byte {
	return k.key.Serialize()
}

func (k *Secp256k1PrivateKey) Hex() string {
	return "0x" + hex.EncodeToString(k.Bytes())
}

// unused rand import guard
var _ = rand.Reader

func (pk *Secp256k1PublicKey) Bytes() []byte {
	return pk.key.SerializeUncompressed()
}

func (pk *Secp256k1PublicKey) Hex() string {
	return "0x" + hex.EncodeToString(pk.Bytes())
}

func (pk *Secp256k1PublicKey) AuthKey() []byte {
	// SingleKey scheme: sha3_256(0x01 || pubkey_bytes || 0x01)
	// variant 0x01 = Secp256k1 in AnyPublicKey, scheme 0x02 = SingleKey
	h := sha3.New256()
	h.Write([]byte{0x01})       // AnyPublicKey variant for Secp256k1
	h.Write(pk.Bytes())
	h.Write([]byte{0x02})       // SingleKey auth scheme
	return h.Sum(nil)
}

func (pk *Secp256k1PublicKey) Serialize(s *bcs.Serializer) {
	s.SerializeFixedBytes(pk.Bytes())
}

func (sig *Secp256k1Signature) Bytes() []byte {
	return sig.bytes[:]
}

func (sig *Secp256k1Signature) Serialize(s *bcs.Serializer) {
	s.SerializeFixedBytes(sig.bytes[:])
}
