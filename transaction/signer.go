package transaction

import (
	"golang.org/x/crypto/sha3"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/bcs"
)

const rawTransactionSalt = "CEDRA::RawTransaction"

// SignTransaction serializes rawTxn, computes the signing message, signs it,
// and returns the BCS-encoded SignedTransaction bytes ready for submission.
func SignTransaction(rawTxn *RawTransaction, signer account.Account) ([]byte, error) {
	txnS := &bcs.Serializer{}
	rawTxn.Serialize(txnS)
	txnBytes := txnS.ToBytes()

	h := sha3.New256()
	h.Write([]byte(rawTransactionSalt))
	h.Write(txnBytes)
	signingMessage := h.Sum(nil)

	authenticator, err := signer.SignTransaction(signingMessage)
	if err != nil {
		return nil, err
	}

	signed := &bcs.Serializer{}
	signed.SerializeFixedBytes(txnBytes)
	signed.SerializeFixedBytes(authenticator)
	return signed.ToBytes(), nil
}
