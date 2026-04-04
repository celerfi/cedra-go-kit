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

// cedraFAAddress is the on-chain FA metadata address for CEDRA (0x000...000a)
const cedraFAAddress = "0x000000000000000000000000000000000000000000000000000000000000000a"

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
	return a.GetAccountFABalance(ctx, address, cedraFAAddress)
}

// GetAccountFABalance returns the balance of any fungible asset for an account using a view function.
func (a *AccountAPI) GetAccountFABalance(ctx context.Context, address, faMetadataAddress string) (uint64, error) {
	req := types.ViewRequest{
		Function:      "0x1::primary_fungible_store::balance",
		TypeArguments: []string{"0x1::fungible_asset::Metadata"},
		Arguments:     []any{address, faMetadataAddress},
	}
	var raw json.RawMessage
	if err := a.c.Post(ctx, "/view", req, &raw); err != nil {
		return 0, fmt.Errorf("account: get FA balance: %w", err)
	}
	var result []any
	if err := json.Unmarshal(raw, &result); err != nil || len(result) == 0 {
		return 0, fmt.Errorf("account: unexpected view response for FA balance")
	}
	switch v := result[0].(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case float64:
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("account: unexpected type %T for FA balance", result[0])
	}
}
