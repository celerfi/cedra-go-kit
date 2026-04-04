package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
)

func TestGenerateEd25519PrivateKey(t *testing.T) {
	k1, err := GenerateEd25519PrivateKey()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	k2, err := GenerateEd25519PrivateKey()
	if err != nil {
		t.Fatalf("generate second: %v", err)
	}
	if hex.EncodeToString(k1.Bytes()) == hex.EncodeToString(k2.Bytes()) {
		t.Error("two generated keys are identical")
	}
}

func TestEd25519PrivateKeyHexRoundTrip(t *testing.T) {
	k, err := GenerateEd25519PrivateKey()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	h := k.Hex()
	k2, err := NewEd25519PrivateKeyFromHex(h)
	if err != nil {
		t.Fatalf("from hex: %v", err)
	}
	if k.Hex() != k2.Hex() {
		t.Errorf("hex round-trip mismatch: %s vs %s", k.Hex(), k2.Hex())
	}
}

func TestEd25519PrivateKeyFromBytesRoundTrip(t *testing.T) {
	k, _ := GenerateEd25519PrivateKey()
	k2, err := NewEd25519PrivateKeyFromBytes(k.Bytes())
	if err != nil {
		t.Fatalf("from bytes: %v", err)
	}
	if k.Hex() != k2.Hex() {
		t.Error("bytes round-trip mismatch")
	}
}

func TestEd25519SignAndVerify(t *testing.T) {
	k, _ := GenerateEd25519PrivateKey()
	msg := []byte("cedra-go-kit test message")
	sig := k.Sign(msg)

	pub := k.PublicKey()
	if !ed25519.Verify(pub.Bytes(), msg, sig.Bytes()) {
		t.Error("signature verification failed")
	}
}

func TestEd25519SignDifferentMessages(t *testing.T) {
	k, _ := GenerateEd25519PrivateKey()
	sig1 := k.Sign([]byte("message one"))
	sig2 := k.Sign([]byte("message two"))

	if hex.EncodeToString(sig1.Bytes()) == hex.EncodeToString(sig2.Bytes()) {
		t.Error("different messages produced same signature")
	}
}

func TestEd25519PublicKeyLength(t *testing.T) {
	k, _ := GenerateEd25519PrivateKey()
	pub := k.PublicKey()
	if len(pub.Bytes()) != 32 {
		t.Errorf("public key should be 32 bytes, got %d", len(pub.Bytes()))
	}
}

func TestEd25519SignatureLength(t *testing.T) {
	k, _ := GenerateEd25519PrivateKey()
	sig := k.Sign([]byte("test"))
	if len(sig.Bytes()) != 64 {
		t.Errorf("signature should be 64 bytes, got %d", len(sig.Bytes()))
	}
}

func TestEd25519AuthKeyLength(t *testing.T) {
	k, _ := GenerateEd25519PrivateKey()
	authKey := k.PublicKey().AuthKey()
	if len(authKey) != 32 {
		t.Errorf("auth key should be 32 bytes, got %d", len(authKey))
	}
}

func TestEd25519AuthKeyDeterministic(t *testing.T) {
	k, _ := GenerateEd25519PrivateKey()
	a1 := k.PublicKey().AuthKey()
	a2 := k.PublicKey().AuthKey()
	if hex.EncodeToString(a1) != hex.EncodeToString(a2) {
		t.Error("auth key is not deterministic")
	}
}

func TestEd25519InvalidPrivateKeyLength(t *testing.T) {
	_, err := NewEd25519PrivateKeyFromBytes([]byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for short private key")
	}
}

func TestEd25519InvalidHex(t *testing.T) {
	_, err := NewEd25519PrivateKeyFromHex("0xzzzz")
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}
