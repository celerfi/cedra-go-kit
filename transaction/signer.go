package transaction

import (
	"golang.org/x/crypto/sha3"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/bcs"
)

const rawTransactionSalt = "CEDRA::RawTransaction"

func SignTransaction(rawTxn *RawTransaction, signer account.Account) ([]byte, error) {
	txnBytes := serializeRawTxn(rawTxn)
	prefix := sha3.Sum256([]byte(rawTransactionSalt))
	signingMessage := append(prefix[:], txnBytes...)

	authenticator, err := signer.SignTransaction(signingMessage)
	if err != nil {
		return nil, err
	}

	signed := &bcs.Serializer{}
	signed.SerializeFixedBytes(txnBytes)
	signed.SerializeFixedBytes(authenticator)
	return signed.ToBytes(), nil
}

func SimulateTransaction(rawTxn *RawTransaction, signer account.Account) ([]byte, error) {
	txnBytes := serializeRawTxn(rawTxn)
	authenticator := zeroedAuthenticator(signer)

	signed := &bcs.Serializer{}
	signed.SerializeFixedBytes(txnBytes)
	signed.SerializeFixedBytes(authenticator)
	return signed.ToBytes(), nil
}

func serializeRawTxn(rawTxn *RawTransaction) []byte {
	s := &bcs.Serializer{}
	rawTxn.Serialize(s)
	return s.ToBytes()
}

func zeroedAuthenticator(signer account.Account) []byte {
	pubKey := signer.PublicKeyBytes()
	s := &bcs.Serializer{}
	s.SerializeULEB128(0)
	s.SerializeBytes(pubKey)
	s.SerializeBytes(make([]byte, 64))
	return s.ToBytes()
}
