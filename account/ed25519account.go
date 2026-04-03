package account

import (
	"github.com/celerfi/cedra-go-kit/bcs"
	"github.com/celerfi/cedra-go-kit/crypto"
)

type Ed25519Account struct {
	privateKey *crypto.Ed25519PrivateKey
	address    AccountAddress
}

func GenerateEd25519Account() (*Ed25519Account, error) {
	priv, err := crypto.GenerateEd25519PrivateKey()
	if err != nil {
		return nil, err
	}
	return NewEd25519AccountFromPrivateKey(priv), nil
}

func NewEd25519AccountFromPrivateKey(priv *crypto.Ed25519PrivateKey) *Ed25519Account {
	authKey := priv.PublicKey().AuthKey()
	var addr AccountAddress
	copy(addr[:], authKey)
	return &Ed25519Account{privateKey: priv, address: addr}
}

func NewEd25519AccountFromHex(hexKey string) (*Ed25519Account, error) {
	priv, err := crypto.NewEd25519PrivateKeyFromHex(hexKey)
	if err != nil {
		return nil, err
	}
	return NewEd25519AccountFromPrivateKey(priv), nil
}

func (a *Ed25519Account) Address() AccountAddress {
	return a.address
}

func (a *Ed25519Account) AuthKey() []byte {
	return a.privateKey.PublicKey().AuthKey()
}

func (a *Ed25519Account) PublicKeyBytes() []byte {
	return a.privateKey.PublicKey().Bytes()
}

func (a *Ed25519Account) PrivateKeyBytes() []byte {
	return a.privateKey.Bytes()
}

func (a *Ed25519Account) PrivateKeyHex() string {
	return a.privateKey.Hex()
}

// SignTransaction returns a BCS-encoded Ed25519 AccountAuthenticator.
// variant 0 = Ed25519
func (a *Ed25519Account) SignTransaction(signingMessage []byte) ([]byte, error) {
	sig := a.privateKey.Sign(signingMessage)
	pub := a.privateKey.PublicKey()

	s := &bcs.Serializer{}
	s.SerializeULEB128(0) // AccountAuthenticatorEd25519 variant
	pub.Serialize(s)
	sig.Serialize(s)
	return s.ToBytes(), nil
}
