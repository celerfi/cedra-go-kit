package account

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestAccountAddressFromHex(t *testing.T) {
	cases := []struct {
		input    string
		wantLen  int
		wantFail bool
	}{
		{"0x1", 32, false},
		{"0x0000000000000000000000000000000000000000000000000000000000000001", 32, false},
		{"0xzzzz", 0, true},
		{"1", 32, false},
	}
	for _, c := range cases {
		addr, err := AccountAddressFromHex(c.input)
		if c.wantFail {
			if err == nil {
				t.Errorf("input %q: expected error", c.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", c.input, err)
			continue
		}
		if len(addr.Bytes()) != c.wantLen {
			t.Errorf("input %q: want len %d got %d", c.input, c.wantLen, len(addr.Bytes()))
		}
	}
}

func TestAccountAddressHexFormat(t *testing.T) {
	addr, _ := AccountAddressFromHex("0x1")
	h := addr.Hex()
	if !strings.HasPrefix(h, "0x") {
		t.Errorf("address hex should start with 0x, got %s", h)
	}
	if len(h) != 66 { // "0x" + 64 hex chars
		t.Errorf("address hex should be 66 chars, got %d: %s", len(h), h)
	}
}

func TestAccountAddressShortPadding(t *testing.T) {
	addr, _ := AccountAddressFromHex("0x1")
	b := addr.Bytes()
	for i := 0; i < 31; i++ {
		if b[i] != 0 {
			t.Errorf("byte %d should be 0, got %d", i, b[i])
		}
	}
	if b[31] != 1 {
		t.Errorf("last byte should be 1, got %d", b[31])
	}
}

func TestAccountAddressRoundTrip(t *testing.T) {
	original := "0x" + strings.Repeat("ab", 32)
	addr, err := AccountAddressFromHex(original)
	if err != nil {
		t.Fatalf("from hex: %v", err)
	}
	if addr.Hex() != original {
		t.Errorf("round-trip: want %s got %s", original, addr.Hex())
	}
}

func TestNewAccountAddressTooLong(t *testing.T) {
	_, err := NewAccountAddress(make([]byte, 33))
	if err == nil {
		t.Error("expected error for 33-byte address")
	}
}

func TestGenerateEd25519Account(t *testing.T) {
	a1, err := GenerateEd25519Account()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	a2, err := GenerateEd25519Account()
	if err != nil {
		t.Fatalf("generate second: %v", err)
	}
	if a1.Address().Hex() == a2.Address().Hex() {
		t.Error("two generated accounts have same address")
	}
}

func TestEd25519AccountAddressLength(t *testing.T) {
	a, _ := GenerateEd25519Account()
	if len(a.Address().Bytes()) != 32 {
		t.Errorf("address should be 32 bytes, got %d", len(a.Address().Bytes()))
	}
}

func TestEd25519AccountAuthKeyMatchesAddress(t *testing.T) {
	a, _ := GenerateEd25519Account()
	if hex.EncodeToString(a.AuthKey()) != hex.EncodeToString(a.Address().Bytes()) {
		t.Error("auth key should match address for freshly generated account")
	}
}

func TestEd25519AccountSignTransaction(t *testing.T) {
	a, _ := GenerateEd25519Account()
	msg := []byte("mock signing message 32 bytes!!!")
	authBytes, err := a.SignTransaction(msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if len(authBytes) == 0 {
		t.Error("authenticator bytes should not be empty")
	}
}

func TestEd25519AccountFromHex(t *testing.T) {
	a, _ := GenerateEd25519Account()
	hex := a.PrivateKeyHex()
	a2, err := NewEd25519AccountFromHex(hex)
	if err != nil {
		t.Fatalf("from hex: %v", err)
	}
	if a.Address().Hex() != a2.Address().Hex() {
		t.Errorf("address mismatch after hex round-trip: %s vs %s", a.Address().Hex(), a2.Address().Hex())
	}
}

func TestGenerateSingleKeyAccount(t *testing.T) {
	a1, err := GenerateSingleKeyAccount()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	a2, err := GenerateSingleKeyAccount()
	if err != nil {
		t.Fatalf("generate second: %v", err)
	}
	if a1.Address().Hex() == a2.Address().Hex() {
		t.Error("two generated single-key accounts have same address")
	}
}

func TestSingleKeyAccountSignTransaction(t *testing.T) {
	a, _ := GenerateSingleKeyAccount()
	msg := []byte("mock signing message 32 bytes!!!")
	authBytes, err := a.SignTransaction(msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if len(authBytes) == 0 {
		t.Error("authenticator bytes should not be empty")
	}
}
