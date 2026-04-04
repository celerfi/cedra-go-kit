package transaction

import (
	"golang.org/x/crypto/sha3"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/bcs"
)

const rawTransactionSalt = "CEDRA::RawTransaction"

// SignTransaction serializes rawTxn, computes the signing message, signs it,
// and returns the BCS-encoded SignedTransaction bytes ready for submission.
//
// Signing message = sha3_256(salt) || bcs(rawTxn)
// The prefix is the hash of the domain separator only; the raw txn bytes are appended unmodified.
func SignTransaction(rawTxn *RawTransaction, signer account.Account) ([]byte, error) {
	txnS := &bcs.Serializer{}
	rawTxn.Serialize(txnS)
	txnBytes := txnS.ToBytes()

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
