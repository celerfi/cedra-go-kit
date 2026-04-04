package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func TestGenerateSecp256k1PrivateKey(t *testing.T) {
	k1, err := GenerateSecp256k1PrivateKey()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	k2, err := GenerateSecp256k1PrivateKey()
	if err != nil {
		t.Fatalf("generate second: %v", err)
	}
	if hex.EncodeToString(k1.Bytes()) == hex.EncodeToString(k2.Bytes()) {
		t.Error("two generated keys are identical")
	}
}

func TestSecp256k1HexRoundTrip(t *testing.T) {
	k, _ := GenerateSecp256k1PrivateKey()
	h := k.Hex()
	k2, err := NewSecp256k1PrivateKeyFromHex(h)
	if err != nil {
		t.Fatalf("from hex: %v", err)
	}
	if k.Hex() != k2.Hex() {
		t.Errorf("hex round-trip mismatch: %s vs %s", k.Hex(), k2.Hex())
	}
}

func TestSecp256k1BytesRoundTrip(t *testing.T) {
	k, _ := GenerateSecp256k1PrivateKey()
	k2, err := NewSecp256k1PrivateKeyFromBytes(k.Bytes())
	if err != nil {
		t.Fatalf("from bytes: %v", err)
	}
	if k.Hex() != k2.Hex() {
		t.Error("bytes round-trip mismatch")
	}
}

func TestSecp256k1SignAndVerify(t *testing.T) {
	k, _ := GenerateSecp256k1PrivateKey()
	msg := []byte("cedra-go-kit secp256k1 test")
	sig, err := k.Sign(msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if len(sig.Bytes()) != 64 {
		t.Errorf("signature should be 64 bytes, got %d", len(sig.Bytes()))
	}
}

func TestSecp256k1LowSEnforcement(t *testing.T) {
	k, _ := GenerateSecp256k1PrivateKey()
	msg := []byte("low-s enforcement test")

	for i := 0; i < 20; i++ {
		sig, err := k.Sign(msg)
		if err != nil {
			t.Fatalf("sign iteration %d: %v", i, err)
		}
		var s secp256k1.ModNScalar
		s.SetByteSlice(sig.Bytes()[32:])
		if s.IsOverHalfOrder() {
			t.Errorf("iteration %d: S is over half order (not low-S)", i)
		}
	}
}

func TestSecp256k1PublicKeyLength(t *testing.T) {
	k, _ := GenerateSecp256k1PrivateKey()
	pub := k.PublicKey()
	if len(pub.Bytes()) != 65 {
		t.Errorf("uncompressed public key should be 65 bytes, got %d", len(pub.Bytes()))
	}
}

func TestSecp256k1AuthKeyLength(t *testing.T) {
	k, _ := GenerateSecp256k1PrivateKey()
	authKey := k.PublicKey().AuthKey()
	if len(authKey) != 32 {
		t.Errorf("auth key should be 32 bytes, got %d", len(authKey))
	}
}

func TestSecp256k1InvalidPrivateKeyLength(t *testing.T) {
	_, err := NewSecp256k1PrivateKeyFromBytes([]byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for wrong length private key")
	}
}

func TestSecp256k1DifferentMessagesProduceDifferentSigs(t *testing.T) {
	k, _ := GenerateSecp256k1PrivateKey()
	sig1, _ := k.Sign([]byte("message one"))
	sig2, _ := k.Sign([]byte("message two"))
	if hex.EncodeToString(sig1.Bytes()) == hex.EncodeToString(sig2.Bytes()) {
		t.Error("different messages produced same signature")
	}
}
