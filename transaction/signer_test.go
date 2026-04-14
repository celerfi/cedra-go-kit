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

func TestFeePayerSigningMessage(t *testing.T) {
	_, feePayerRawTxn, _, _ := makeFixedFeePayerRawTxn(t)

	const wantHex = "dc1cd51bc0103d1c06a90c8ad02db5b5b1c4e4dab7a1f3a96b3a93d7b0cd315901147e4d3a5b10eaed2a93536e284c23096dfcea9ac61f0a8420e5d01fbd8f0ea807000000000000000200000000000000000000000000000000000000000000000000000000000000010d63656472615f6163636f756e74087472616e7366657200022000000000000000000000000000000000000000000000000000000000deadbeef089210000000000000400d030000000000640000000000000000f1536500000000020700000000000000000000000000000000000000000000000000000000000000010a63656472615f636f696e094365647261436f696e00000000000000000000000000000000000000000000000000000000000000000000"
	if got := hex.EncodeToString(FeePayerSigningMessage(feePayerRawTxn)); got != wantHex {
		t.Fatalf("fee-payer signing message mismatch\nwant: %s\ngot:  %s", wantHex, got)
	}
}

func TestSignFeePayerTransactionDeterministic(t *testing.T) {
	_, feePayerRawTxn, sender, feePayer := makeFixedFeePayerRawTxn(t)

	signed1, err := SignFeePayerTransaction(feePayerRawTxn, sender, feePayer)
	if err != nil {
		t.Fatalf("sign fee-payer transaction: %v", err)
	}
	signed2, err := SignFeePayerTransaction(feePayerRawTxn, sender, feePayer)
	if err != nil {
		t.Fatalf("sign fee-payer transaction again: %v", err)
	}

	const wantHex = "147e4d3a5b10eaed2a93536e284c23096dfcea9ac61f0a8420e5d01fbd8f0ea807000000000000000200000000000000000000000000000000000000000000000000000000000000010d63656472615f6163636f756e74087472616e7366657200022000000000000000000000000000000000000000000000000000000000deadbeef089210000000000000400d030000000000640000000000000000f1536500000000020700000000000000000000000000000000000000000000000000000000000000010a63656472615f636f696e094365647261436f696e00030020d04ab232742bb4ab3a1368bd4615e4e6d0224ab71a016baf8520a332c977873740ff37089000ccc20d6169f8500e4ea80c74e95b117b338ee6c6e8f17199936a9e75aa15d4adead985cc24e996dfc26b0067aee82c5c719ebc8d64d796b795c7070000a32657fd60acb0433491a33d84823c04722ae76639b272873cc27d015232904e0020a09aa5f47a6759802ff955f8dc2d2a14a5c99d23be97f864127ff9383455a4f0409c4cb89287438ee0450dd0917cf80b4069d35d075d6e9b0c02c35883c39e616c130fcf23b5790bcbf5b0b24b17012f2a36deeb63f75500840ea7ba14636b1601"

	got1 := hex.EncodeToString(signed1)
	got2 := hex.EncodeToString(signed2)
	if got1 != wantHex {
		t.Fatalf("fee-payer signed txn mismatch\nwant: %s\ngot:  %s", wantHex, got1)
	}
	if got2 != wantHex {
		t.Fatalf("fee-payer signed txn changed across runs\nwant: %s\ngot:  %s", wantHex, got2)
	}
}

func TestAssembleFeePayerSignedTransactionUsesExplicitFeePayerAddress(t *testing.T) {
	_, feePayerRawTxn, sender, feePayer := makeFixedFeePayerRawTxn(t)

	senderAuthenticator, err := SignFeePayerTransactionSenderAuthenticator(feePayerRawTxn, sender)
	if err != nil {
		t.Fatalf("sign sender authenticator: %v", err)
	}
	feePayerAuthenticator, err := SignFeePayerTransactionFeePayerAuthenticator(feePayerRawTxn, feePayer)
	if err != nil {
		t.Fatalf("sign fee payer authenticator: %v", err)
	}

	signed, err := AssembleFeePayerSignedTransaction(feePayerRawTxn, senderAuthenticator, feePayer.Address(), feePayerAuthenticator)
	if err != nil {
		t.Fatalf("assemble fee payer transaction: %v", err)
	}

	authenticator := signed[len(serializeRawTxn(feePayerRawTxn.RawTransaction)):]
	if len(authenticator) == 0 || authenticator[0] != 0x03 {
		t.Fatalf("expected fee-payer authenticator variant 3, got %x", authenticator)
	}
}

func TestSimulateTransactionSingleKeyAuthenticatorShape(t *testing.T) {
	signer, err := account.NewSingleKeyAccountFromHex("0x0101010101010101010101010101010101010101010101010101010101010101")
	if err != nil {
		t.Fatalf("single key account: %v", err)
	}
	rawTxn := makeTestRawTxn(t)

	simulated, err := SimulateTransaction(rawTxn, signer)
	if err != nil {
		t.Fatalf("simulate transaction: %v", err)
	}

	authenticator := simulated[len(serializeRawTxn(rawTxn)):]
	if len(authenticator) == 0 || authenticator[0] != 0x02 {
		t.Fatalf("expected single-key authenticator variant 2, got %x", authenticator)
	}
}
