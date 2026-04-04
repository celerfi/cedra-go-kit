package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/celerfi/cedra-go-kit/client"
	"github.com/celerfi/cedra-go-kit/types"
)

type ANSAPI struct {
	c       *client.Client
	general *GeneralAPI
}

func newANSAPI(c *client.Client, g *GeneralAPI) *ANSAPI {
	return &ANSAPI{c: c, general: g}
}

// GetAddressFromName resolves a .cedra name to an account address.
// name should be in format "alice" or "alice.cedra"
func (a *ANSAPI) GetAddressFromName(ctx context.Context, name string) (string, error) {
	name = strings.TrimSuffix(name, ".cedra")
	req := types.ViewRequest{
		Function:      "0x867ed1f6bf916171b1de3ee92849b8978b7d1b9e0a8cc982a3d19d535dfd9c0c::domains::get_name_address",
		TypeArguments: []string{},
		Arguments:     []any{name, ""},
	}
	raw, err := a.general.View(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ans: resolve name %q: %w", name, err)
	}
	var result []json.RawMessage
	if err := json.Unmarshal(raw, &result); err != nil || len(result) == 0 {
		return "", fmt.Errorf("ans: unexpected response for name %q", name)
	}
	var addr string
	if err := json.Unmarshal(result[0], &addr); err != nil {
		return "", fmt.Errorf("ans: parse address for name %q: %w", name, err)
	}
	return addr, nil
}

// GetNameFromAddress performs reverse lookup: address → primary .cedra name
func (a *ANSAPI) GetNameFromAddress(ctx context.Context, address string) (string, error) {
	req := types.ViewRequest{
		Function:      "0x867ed1f6bf916171b1de3ee92849b8978b7d1b9e0a8cc982a3d19d535dfd9c0c::domains::get_reverse_lookup",
		TypeArguments: []string{},
		Arguments:     []any{address},
	}
	raw, err := a.general.View(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ans: reverse lookup %q: %w", address, err)
	}
	var result []json.RawMessage
	if err := json.Unmarshal(raw, &result); err != nil || len(result) == 0 {
		return "", fmt.Errorf("ans: unexpected response for address %q", address)
	}
	var name string
	if err := json.Unmarshal(result[0], &name); err != nil {
		return "", fmt.Errorf("ans: parse name for address %q: %w", address, err)
	}
	if name != "" {
		name = name + ".cedra"
	}
	return name, nil
}
