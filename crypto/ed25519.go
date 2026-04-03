package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/celerfi/cedra-go-kit/bcs"
	"golang.org/x/crypto/sha3"
)

type Ed25519PrivateKey struct {
	key ed25519.PrivateKey
}

type Ed25519PublicKey struct {
	key ed25519.PublicKey
}

type Ed25519Signature struct {
	bytes [64]byte
}

func GenerateEd25519PrivateKey() (*Ed25519PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Ed25519PrivateKey{key: priv}, nil
}

func NewEd25519PrivateKeyFromBytes(b []byte) (*Ed25519PrivateKey, error) {
	if len(b) != ed25519.SeedSize && len(b) != ed25519.PrivateKeySize {
		return nil, errors.New("ed25519: invalid private key length")
	}
	if len(b) == ed25519.SeedSize {
		return &Ed25519PrivateKey{key: ed25519.NewKeyFromSeed(b)}, nil
	}
	return &Ed25519PrivateKey{key: ed25519.PrivateKey(b)}, nil
}

func NewEd25519PrivateKeyFromHex(s string) (*Ed25519PrivateKey, error) {
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return NewEd25519PrivateKeyFromBytes(b)
}

func (k *Ed25519PrivateKey) PublicKey() *Ed25519PublicKey {
	return &Ed25519PublicKey{key: k.key.Public().(ed25519.PublicKey)}
}

func (k *Ed25519PrivateKey) Sign(message []byte) *Ed25519Signature {
	sig := ed25519.Sign(k.key, message)
	var out Ed25519Signature
	copy(out.bytes[:], sig)
	return &out
}

func (k *Ed25519PrivateKey) Seed() []byte {
	return k.key.Seed()
}

func (k *Ed25519PrivateKey) Bytes() []byte {
	return k.key.Seed()
}

func (k *Ed25519PrivateKey) Hex() string {
	return "0x" + hex.EncodeToString(k.Bytes())
}

func (pk *Ed25519PublicKey) Bytes() []byte {
	return []byte(pk.key)
}

func (pk *Ed25519PublicKey) Hex() string {
	return "0x" + hex.EncodeToString(pk.Bytes())
}

func (pk *Ed25519PublicKey) AuthKey() []byte {
	h := sha3.New256()
	h.Write(pk.Bytes())
	h.Write([]byte{0x00})
	return h.Sum(nil)
}

func (pk *Ed25519PublicKey) Serialize(s *bcs.Serializer) {
	s.SerializeFixedBytes(pk.Bytes())
}

func (sig *Ed25519Signature) Bytes() []byte {
	return sig.bytes[:]
}

func (sig *Ed25519Signature) Serialize(s *bcs.Serializer) {
	s.SerializeFixedBytes(sig.bytes[:])
}
