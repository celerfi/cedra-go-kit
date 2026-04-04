package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/celerfi/cedra-go-kit/client"
)

type FaucetAPI struct {
	c *client.Client
}

func newFaucetAPI(c *client.Client) *FaucetAPI {
	return &FaucetAPI{c: c}
}

func (f *FaucetAPI) FundAccount(ctx context.Context, address string, amount uint64) error {
	if f.c.FaucetURL() == "" {
		return errors.New("faucet: no faucet URL configured for this network")
	}
	body := map[string]any{
		"address": address,
		"amount":  amount,
	}
	var result any
	if err := f.c.PostFaucet(ctx, "/mint", body, &result); err != nil {
		return fmt.Errorf("faucet: fund account: %w", err)
	}
	return nil
}
