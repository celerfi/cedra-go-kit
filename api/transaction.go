package api

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/client"
	"github.com/celerfi/cedra-go-kit/transaction"
	"github.com/celerfi/cedra-go-kit/types"
)

type TransactionAPI struct {
	c       *client.Client
	builder *transaction.Builder
}

func newTransactionAPI(c *client.Client) *TransactionAPI {
	return &TransactionAPI{c: c, builder: transaction.NewBuilder(c)}
}

func (t *TransactionAPI) GetTransactionByHash(ctx context.Context, hash string) (*types.CommittedTransaction, error) {
	var txn types.CommittedTransaction
	if err := t.c.Get(ctx, "/transactions/by_hash/"+hash, nil, &txn); err != nil {
		return nil, err
	}
	return &txn, nil
}

func (t *TransactionAPI) GetTransactionByVersion(ctx context.Context, version uint64) (*types.CommittedTransaction, error) {
	var txn types.CommittedTransaction
	if err := t.c.Get(ctx, "/transactions/by_version/"+strconv.FormatUint(version, 10), nil, &txn); err != nil {
		return nil, err
	}
	return &txn, nil
}

func (t *TransactionAPI) GetTransactions(ctx context.Context, limit, start *uint64) ([]types.CommittedTransaction, error) {
	params := url.Values{}
	if limit != nil {
		params.Set("limit", strconv.FormatUint(*limit, 10))
	}
	if start != nil {
		params.Set("start", strconv.FormatUint(*start, 10))
	}
	var txns []types.CommittedTransaction
	if err := t.c.Get(ctx, "/transactions", params, &txns); err != nil {
		return nil, err
	}
	return txns, nil
}

func (t *TransactionAPI) SubmitTransaction(ctx context.Context, signedTxnBCS []byte) (*types.PendingTransaction, error) {
	var pending types.PendingTransaction
	if err := t.c.PostBCS(ctx, "/transactions", signedTxnBCS, &pending); err != nil {
		return nil, err
	}
	return &pending, nil
}

func (t *TransactionAPI) BuildTransaction(ctx context.Context, sender account.AccountAddress, opts transaction.BuildOptions) (*transaction.RawTransaction, error) {
	return t.builder.Build(ctx, sender, opts)
}

func (t *TransactionAPI) SimulateTransaction(ctx context.Context, rawTxn *transaction.RawTransaction, signer account.Account) ([]types.CommittedTransaction, error) {
	signedBytes, err := transaction.SimulateTransaction(rawTxn, signer)
	if err != nil {
		return nil, err
	}
	var results []types.CommittedTransaction
	if err := t.c.PostBCS(ctx, "/transactions/simulate", signedBytes, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (t *TransactionAPI) IsTransactionPending(ctx context.Context, hash string) (bool, error) {
	txn, err := t.GetTransactionByHash(ctx, hash)
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return true, nil
		}
		return false, err
	}
	return txn.Type == "pending_transaction", nil
}

func (t *TransactionAPI) WaitForTransaction(ctx context.Context, hash string) (*types.CommittedTransaction, error) {
	deadline := time.Now().Add(time.Duration(client.DefaultTxnTimeoutSecs) * time.Second)
	for time.Now().Before(deadline) {
		txn, err := t.GetTransactionByHash(ctx, hash)
		if err != nil {
			var apiErr *client.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return nil, err
		}
		if txn.Type != "pending_transaction" {
			if !txn.Success {
				return txn, fmt.Errorf("transaction failed: %s", txn.VmStatus)
			}
			return txn, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil, fmt.Errorf("transaction %s timed out after %d seconds", hash, client.DefaultTxnTimeoutSecs)
}
