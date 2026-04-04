package transaction

import (
	"encoding/hex"
	"testing"

	"github.com/celerfi/cedra-go-kit/account"
)

func makeTestRawTxn(t *testing.T) *RawTransaction {
	t.Helper()
	sender, _ := account.AccountAddressFromHex("0xdeadbeef")
	moduleAddr, _ := account.AccountAddressFromHex("0x1")
	return &RawTransaction{
		Sender:         sender,
		SequenceNumber: 0,
		Payload: &EntryFunction{
			Module:   ModuleID{Address: moduleAddr, Name: "cedra_account"},
			Function: "transfer",
			TypeArgs: []TypeTag{},
			Args: [][]byte{
				SerializeAddressArg(sender),
				SerializeU64Arg(1_000_000),
			},
		},
		MaxGasAmount:            200_000,
		GasUnitPrice:            100,
		ExpirationTimestampSecs: 9999999999,
		ChainID:                 2,
	}
}

func TestSignTransactionEd25519NonEmpty(t *testing.T) {
	signer, err := account.GenerateEd25519Account()
	if err != nil {
		t.Fatalf("generate account: %v", err)
	}
	rawTxn := makeTestRawTxn(t)
	signed, err := SignTransaction(rawTxn, signer)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if len(signed) == 0 {
		t.Error("signed transaction bytes should not be empty")
	}
}

func TestSignTransactionSingleKeyNonEmpty(t *testing.T) {
	signer, err := account.GenerateSingleKeyAccount()
	if err != nil {
		t.Fatalf("generate account: %v", err)
	}
	rawTxn := makeTestRawTxn(t)
	signed, err := SignTransaction(rawTxn, signer)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if len(signed) == 0 {
		t.Error("signed transaction bytes should not be empty")
	}
}

func TestSignTransactionDeterministicEd25519(t *testing.T) {
	signer, _ := account.GenerateEd25519Account()
	rawTxn := makeTestRawTxn(t)

	s1, _ := SignTransaction(rawTxn, signer)
	s2, _ := SignTransaction(rawTxn, signer)

	if hex.EncodeToString(s1) != hex.EncodeToString(s2) {
		t.Error("Ed25519 signed transaction is not deterministic for same key+txn")
	}
}

func TestSignTransactionDifferentSigners(t *testing.T) {
	s1, _ := account.GenerateEd25519Account()
	s2, _ := account.GenerateEd25519Account()
	rawTxn := makeTestRawTxn(t)

	b1, _ := SignTransaction(rawTxn, s1)
	b2, _ := SignTransaction(rawTxn, s2)

	if hex.EncodeToString(b1) == hex.EncodeToString(b2) {
		t.Error("different signers produced identical signed transactions")
	}
}
