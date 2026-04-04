package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/celerfi/cedra-go-kit/client"
)

type FaucetAPI struct {
	c       *client.Client
	txnAPI  *TransactionAPI
}

func newFaucetAPI(c *client.Client, txnAPI *TransactionAPI) *FaucetAPI {
	return &FaucetAPI{c: c, txnAPI: txnAPI}
}

func (f *FaucetAPI) FundAccount(ctx context.Context, address string, amount uint64) error {
	if f.c.FaucetURL() == "" {
		return errors.New("faucet: no faucet URL configured for this network")
	}
	body := map[string]any{
		"address": address,
		"amount":  amount,
	}
	var result struct {
		TxnHashes []string `json:"txn_hashes"`
	}
	if err := f.c.PostFaucet(ctx, "/fund", body, &result); err != nil {
		return fmt.Errorf("faucet: fund account: %w", err)
	}
	if len(result.TxnHashes) == 0 {
		return errors.New("faucet: no transaction hash returned")
	}
	deadline := time.Now().Add(30 * time.Second)
	for _, hash := range result.TxnHashes {
		waitCtx := ctx
		var cancel context.CancelFunc
		if dl, ok := ctx.Deadline(); !ok || dl.After(deadline) {
			waitCtx, cancel = context.WithDeadline(ctx, deadline)
			defer cancel()
		}
		txn, err := f.txnAPI.WaitForTransaction(waitCtx, hash)
		if err != nil {
			return fmt.Errorf("faucet: wait for txn %s: %w", hash, err)
		}
		if !txn.Success {
			return fmt.Errorf("faucet: txn %s failed: %s", hash, txn.VmStatus)
		}
	}
	return nil
}

// FundAccountNoWait funds an account and returns without waiting for the transaction.
func (f *FaucetAPI) FundAccountNoWait(ctx context.Context, address string, amount uint64) ([]string, error) {
	if f.c.FaucetURL() == "" {
		return nil, errors.New("faucet: no faucet URL configured for this network")
	}
	body := map[string]any{
		"address": address,
		"amount":  amount,
	}
	var result struct {
		TxnHashes []string `json:"txn_hashes"`
	}
	if err := f.c.PostFaucet(ctx, "/fund", body, &result); err != nil {
		return nil, fmt.Errorf("faucet: fund account: %w", err)
	}
	return result.TxnHashes, nil
}
