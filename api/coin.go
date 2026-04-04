package api

import (
	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/transaction"
	"github.com/celerfi/cedra-go-kit/types"
)

type CoinAPI struct{}

func newCoinAPI() *CoinAPI { return &CoinAPI{} }

func (c *CoinAPI) TransferCEDRA(recipient account.AccountAddress, amount uint64, opts *types.TransactionOptions) transaction.BuildOptions {
	return transaction.BuildOptions{
		Function: "0x1::cedra_account::transfer",
		TypeArgs: []string{},
		Args: [][]byte{
			transaction.SerializeAddressArg(recipient),
			transaction.SerializeU64Arg(amount),
		},
		Options: opts,
	}
}

func (c *CoinAPI) TransferCoin(coinType string, recipient account.AccountAddress, amount uint64, opts *types.TransactionOptions) transaction.BuildOptions {
	return transaction.BuildOptions{
		Function: "0x1::coin::transfer",
		TypeArgs: []string{coinType},
		Args: [][]byte{
			transaction.SerializeAddressArg(recipient),
			transaction.SerializeU64Arg(amount),
		},
		Options: opts,
	}
}
