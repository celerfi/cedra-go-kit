package api

import (
	"context"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/client"
	"github.com/celerfi/cedra-go-kit/transaction"
	"github.com/celerfi/cedra-go-kit/types"
)

type Cedra struct {
	cfg         client.Config
	httpClient  *client.Client
	General     *GeneralAPI
	Account     *AccountAPI
	Transaction *TransactionAPI
	Event       *EventAPI
	Coin        *CoinAPI
	Faucet      *FaucetAPI
	ANS         *ANSAPI
}

func NewCedra(network client.Network) *Cedra {
	return NewCedraWithConfig(client.DefaultConfig(network))
}

func NewCedraWithConfig(cfg client.Config) *Cedra {
	c := client.NewClient(cfg)
	g := newGeneralAPI(c)
	txnAPI := newTransactionAPI(c)
	return &Cedra{
		cfg:         cfg,
		httpClient:  c,
		General:     g,
		Account:     newAccountAPI(c),
		Transaction: txnAPI,
		Event:       newEventAPI(c),
		Coin:        newCoinAPI(),
		Faucet:      newFaucetAPI(c, txnAPI),
		ANS:         newANSAPI(c, g),
	}
}

func (ce *Cedra) SignAndSubmitTransaction(ctx context.Context, signer account.Account, opts transaction.BuildOptions) (*types.CommittedTransaction, error) {
	rawTxn, err := ce.Transaction.BuildTransaction(ctx, signer.Address(), opts)
	if err != nil {
		return nil, err
	}
	signedBytes, err := transaction.SignTransaction(rawTxn, signer)
	if err != nil {
		return nil, err
	}
	pending, err := ce.Transaction.SubmitTransaction(ctx, signedBytes)
	if err != nil {
		return nil, err
	}
	return ce.Transaction.WaitForTransaction(ctx, pending.Hash)
}

func (ce *Cedra) WaitForTransaction(ctx context.Context, hash string) (*types.CommittedTransaction, error) {
	return ce.Transaction.WaitForTransaction(ctx, hash)
}
