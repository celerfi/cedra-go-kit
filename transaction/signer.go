package transaction

import (
	"fmt"

	"golang.org/x/crypto/sha3"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/bcs"
)

const rawTransactionSalt = "CEDRA::RawTransaction"
const rawTransactionWithDataSalt = "CEDRA::RawTransactionWithData"
const feePayerTransactionAuthenticatorVariant = 3

func SignTransaction(rawTxn *RawTransaction, signer account.Account) ([]byte, error) {
	authenticator, err := SignTransactionAuthenticator(rawTxn, signer)
	if err != nil {
		return nil, err
	}
	return AssembleSignedTransaction(rawTxn, authenticator), nil
}

func SignFeePayerTransaction(rawTxn *FeePayerRawTransaction, sender account.Account, feePayer account.Account) ([]byte, error) {
	senderAuthenticator, err := SignFeePayerTransactionSenderAuthenticator(rawTxn, sender)
	if err != nil {
		return nil, err
	}
	feePayerAuthenticator, err := SignFeePayerTransactionFeePayerAuthenticator(rawTxn, feePayer)
	if err != nil {
		return nil, err
	}
	return AssembleFeePayerSignedTransaction(rawTxn, senderAuthenticator, feePayer.Address(), feePayerAuthenticator)
}

func SimulateTransaction(rawTxn *RawTransaction, signer account.Account) ([]byte, error) {
	return AssembleSignedTransaction(rawTxn, zeroedAuthenticator(signer)), nil
}

func SimulateFeePayerTransaction(rawTxn *FeePayerRawTransaction, sender account.Account, feePayer account.Account) ([]byte, error) {
	return AssembleFeePayerSignedTransaction(rawTxn, zeroedAuthenticator(sender), feePayer.Address(), zeroedAuthenticator(feePayer))
}

func SigningMessage(rawTxn *RawTransaction) []byte {
	return signingMessage(rawTransactionSalt, serializeRawTxn(rawTxn))
}

func FeePayerSigningMessage(rawTxn *FeePayerRawTransaction) []byte {
	return signingMessage(rawTransactionWithDataSalt, serializeFeePayerRawTxn(rawTxn))
}

func SignTransactionAuthenticator(rawTxn *RawTransaction, signer account.Account) ([]byte, error) {
	return signer.SignTransaction(SigningMessage(rawTxn))
}

func SignFeePayerTransactionSenderAuthenticator(rawTxn *FeePayerRawTransaction, signer account.Account) ([]byte, error) {
	return signer.SignTransaction(FeePayerSigningMessage(rawTxn))
}

func SignFeePayerTransactionFeePayerAuthenticator(rawTxn *FeePayerRawTransaction, signer account.Account) ([]byte, error) {
	return signer.SignTransaction(FeePayerSigningMessage(feePayerSigningTransaction(rawTxn, signer.Address())))
}

func AssembleSignedTransaction(rawTxn *RawTransaction, authenticator []byte) []byte {
	return assembleSignedTransaction(serializeRawTxn(rawTxn), authenticator)
}

func AssembleFeePayerSignedTransaction(rawTxn *FeePayerRawTransaction, senderAuthenticator []byte, feePayerAddress account.AccountAddress, feePayerAuthenticator []byte) ([]byte, error) {
	if rawTxn == nil || rawTxn.RawTransaction == nil {
		return nil, fmt.Errorf("transaction: raw transaction is required")
	}

	authenticator := &bcs.Serializer{}
	authenticator.SerializeULEB128(feePayerTransactionAuthenticatorVariant)
	authenticator.SerializeFixedBytes(senderAuthenticator)
	authenticator.SerializeULEB128(0)
	authenticator.SerializeULEB128(0)
	feePayerAddress.Serialize(authenticator)
	authenticator.SerializeFixedBytes(feePayerAuthenticator)

	return assembleSignedTransaction(serializeRawTxn(rawTxn.RawTransaction), authenticator.ToBytes()), nil
}

func serializeRawTxn(rawTxn *RawTransaction) []byte {
	s := &bcs.Serializer{}
	rawTxn.Serialize(s)
	return s.ToBytes()
}

func serializeFeePayerRawTxn(rawTxn *FeePayerRawTransaction) []byte {
	s := &bcs.Serializer{}
	rawTxn.Serialize(s)
	return s.ToBytes()
}

func feePayerSigningTransaction(rawTxn *FeePayerRawTransaction, feePayerAddress account.AccountAddress) *FeePayerRawTransaction {
	return &FeePayerRawTransaction{
		RawTransaction:  rawTxn.RawTransaction,
		FeePayerAddress: feePayerAddress,
	}
}

func zeroedAuthenticator(signer account.Account) []byte {
	pubKey := signer.PublicKeyBytes()
	s := &bcs.Serializer{}
	switch signer.(type) {
	case *account.SingleKeyAccount:
		s.SerializeULEB128(2)
		s.SerializeULEB128(1)
		s.SerializeBytes(pubKey)
		s.SerializeULEB128(1)
		s.SerializeBytes(make([]byte, 64))
	default:
		s.SerializeULEB128(0)
		s.SerializeBytes(pubKey)
		s.SerializeBytes(make([]byte, 64))
	}
	return s.ToBytes()
}

func signingMessage(salt string, txnBytes []byte) []byte {
	prefix := sha3.Sum256([]byte(salt))
	return append(prefix[:], txnBytes...)
}

func assembleSignedTransaction(txnBytes []byte, authenticator []byte) []byte {
	signed := &bcs.Serializer{}
	signed.SerializeFixedBytes(txnBytes)
	signed.SerializeFixedBytes(authenticator)
	return signed.ToBytes()
}
