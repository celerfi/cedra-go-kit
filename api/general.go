package api

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/celerfi/cedra-go-kit/client"
	"github.com/celerfi/cedra-go-kit/types"
)

type GeneralAPI struct {
	c *client.Client
}

func newGeneralAPI(c *client.Client) *GeneralAPI {
	return &GeneralAPI{c: c}
}

func (g *GeneralAPI) GetLedgerInfo(ctx context.Context) (*types.LedgerInfo, error) {
	var info types.LedgerInfo
	if err := g.c.Get(ctx, "/", nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (g *GeneralAPI) GetChainID(ctx context.Context) (uint8, error) {
	info, err := g.GetLedgerInfo(ctx)
	if err != nil {
		return 0, err
	}
	return info.ChainID, nil
}

func (g *GeneralAPI) GetGasEstimation(ctx context.Context) (*types.GasEstimation, error) {
	var est types.GasEstimation
	if err := g.c.Get(ctx, "/estimate_gas_price", nil, &est); err != nil {
		return nil, err
	}
	return &est, nil
}

func (g *GeneralAPI) GetBlockByHeight(ctx context.Context, height uint64, withTxns bool) (*types.Block, error) {
	params := url.Values{}
	params.Set("with_transactions", strconv.FormatBool(withTxns))
	var block types.Block
	if err := g.c.Get(ctx, "/blocks/by_height/"+strconv.FormatUint(height, 10), params, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func (g *GeneralAPI) GetBlockByVersion(ctx context.Context, version uint64, withTxns bool) (*types.Block, error) {
	params := url.Values{}
	params.Set("with_transactions", strconv.FormatBool(withTxns))
	var block types.Block
	if err := g.c.Get(ctx, "/blocks/by_version/"+strconv.FormatUint(version, 10), params, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func (g *GeneralAPI) View(ctx context.Context, req types.ViewRequest) (json.RawMessage, error) {
	var result json.RawMessage
	if err := g.c.Post(ctx, "/view", req, &result); err != nil {
		return nil, err
	}
	return result, nil
}
