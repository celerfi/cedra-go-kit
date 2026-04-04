package account

import (
	"github.com/celerfi/cedra-go-kit/bcs"
	"github.com/celerfi/cedra-go-kit/crypto"
)

type SingleKeyAccount struct {
	privateKey *crypto.Secp256k1PrivateKey
	address    AccountAddress
}

func GenerateSingleKeyAccount() (*SingleKeyAccount, error) {
	priv, err := crypto.GenerateSecp256k1PrivateKey()
	if err != nil {
		return nil, err
	}
	return NewSingleKeyAccountFromPrivateKey(priv), nil
}

func NewSingleKeyAccountFromPrivateKey(priv *crypto.Secp256k1PrivateKey) *SingleKeyAccount {
	authKey := priv.PublicKey().AuthKey()
	var addr AccountAddress
	copy(addr[:], authKey)
	return &SingleKeyAccount{privateKey: priv, address: addr}
}

func NewSingleKeyAccountFromHex(hexKey string) (*SingleKeyAccount, error) {
	priv, err := crypto.NewSecp256k1PrivateKeyFromHex(hexKey)
	if err != nil {
		return nil, err
	}
	return NewSingleKeyAccountFromPrivateKey(priv), nil
}

func (a *SingleKeyAccount) Address() AccountAddress  { return a.address }
func (a *SingleKeyAccount) AuthKey() []byte          { return a.privateKey.PublicKey().AuthKey() }
func (a *SingleKeyAccount) PublicKeyBytes() []byte   { return a.privateKey.PublicKey().Bytes() }
func (a *SingleKeyAccount) PrivateKeyBytes() []byte  { return a.privateKey.Bytes() }
func (a *SingleKeyAccount) PrivateKeyHex() string    { return a.privateKey.Hex() }

func (a *SingleKeyAccount) SignTransaction(signingMessage []byte) ([]byte, error) {
	sig, err := a.privateKey.Sign(signingMessage)
	if err != nil {
		return nil, err
	}
	pub := a.privateKey.PublicKey()
	s := &bcs.Serializer{}
	s.SerializeULEB128(2)
	s.SerializeULEB128(1)
	pub.Serialize(s)
	s.SerializeULEB128(1)
	sig.Serialize(s)
	return s.ToBytes(), nil
}
