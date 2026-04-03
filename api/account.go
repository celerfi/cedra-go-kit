package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/celerfi/cedra-go-kit/client"
	"github.com/celerfi/cedra-go-kit/types"
)

type AccountAPI struct {
	c *client.Client
}

func newAccountAPI(c *client.Client) *AccountAPI {
	return &AccountAPI{c: c}
}

func (a *AccountAPI) GetAccountInfo(ctx context.Context, address string) (*types.AccountData, error) {
	var acct types.AccountData
	if err := a.c.Get(ctx, "/accounts/"+address, nil, &acct); err != nil {
		return nil, err
	}
	return &acct, nil
}

func (a *AccountAPI) GetSequenceNumber(ctx context.Context, address string) (uint64, error) {
	acct, err := a.GetAccountInfo(ctx, address)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(acct.SequenceNumber, 10, 64)
}

func (a *AccountAPI) GetAccountResources(ctx context.Context, address string) ([]types.MoveResource, error) {
	var resources []types.MoveResource
	if err := a.c.Get(ctx, "/accounts/"+address+"/resources", nil, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}

func (a *AccountAPI) GetAccountResource(ctx context.Context, address, resourceType string) (*types.MoveResource, error) {
	var resource types.MoveResource
	if err := a.c.Get(ctx, "/accounts/"+address+"/resource/"+url.PathEscape(resourceType), nil, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

func (a *AccountAPI) GetAccountModules(ctx context.Context, address string) ([]types.MoveModuleBytecode, error) {
	var modules []types.MoveModuleBytecode
	if err := a.c.Get(ctx, "/accounts/"+address+"/modules", nil, &modules); err != nil {
		return nil, err
	}
	return modules, nil
}

func (a *AccountAPI) GetAccountTransactions(ctx context.Context, address string, limit, start *uint64) ([]types.CommittedTransaction, error) {
	params := url.Values{}
	if limit != nil {
		params.Set("limit", strconv.FormatUint(*limit, 10))
	}
	if start != nil {
		params.Set("start", strconv.FormatUint(*start, 10))
	}
	var txns []types.CommittedTransaction
	if err := a.c.Get(ctx, "/accounts/"+address+"/transactions", params, &txns); err != nil {
		return nil, err
	}
	return txns, nil
}

func (a *AccountAPI) GetAccountCEDRABalance(ctx context.Context, address string) (uint64, error) {
	resource, err := a.GetAccountResource(ctx, address, "0x1::coin::CoinStore<0x1::cedra_coin::CedraCoin>")
	if err != nil {
		return 0, err
	}
	var store types.CoinStore
	if err := json.Unmarshal(resource.Data, &store); err != nil {
		return 0, fmt.Errorf("account: parse coin store: %w", err)
	}
	return strconv.ParseUint(store.Coin.Value, 10, 64)
}
