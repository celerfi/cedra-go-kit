package transaction

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/celerfi/cedra-go-kit/account"
	"github.com/celerfi/cedra-go-kit/client"
	"github.com/celerfi/cedra-go-kit/types"
)

type BuildOptions struct {
	Function        string
	TypeArgs        []string
	Args            [][]byte
	WithFeePayer    bool
	FeePayerAddress *account.AccountAddress
	Options         *types.TransactionOptions
}

type Builder struct {
	c *client.Client
}

func NewBuilder(c *client.Client) *Builder {
	return &Builder{c: c}
}

func (b *Builder) Build(ctx context.Context, sender account.AccountAddress, opts BuildOptions) (*RawTransaction, error) {
	if opts.WithFeePayer {
		return nil, fmt.Errorf("builder: with_fee_payer requested; use BuildFeePayer")
	}
	if opts.FeePayerAddress != nil {
		return nil, fmt.Errorf("builder: fee payer address provided; use BuildFeePayer")
	}
	return b.buildRawTransaction(ctx, sender, opts)
}

func (b *Builder) BuildFeePayer(ctx context.Context, sender account.AccountAddress, opts BuildOptions) (*FeePayerRawTransaction, error) {
	rawTxn, err := b.buildRawTransaction(ctx, sender, opts)
	if err != nil {
		return nil, err
	}
	var feePayer account.AccountAddress
	if opts.FeePayerAddress != nil {
		feePayer = *opts.FeePayerAddress
	}
	return rawTxn.WithFeePayer(feePayer), nil
}

func (b *Builder) buildRawTransaction(ctx context.Context, sender account.AccountAddress, opts BuildOptions) (*RawTransaction, error) {
	var seq uint64
	fetchFromNetwork := true

	// Check if a manual Sequence Number was pushed via options
	// This skips the network call, preventing mempool nonce collisions on highly concurrent events
	if opts.Options != nil && opts.Options.SequenceNumber != nil {
		seq = *opts.Options.SequenceNumber
		fetchFromNetwork = false
	}

	if fetchFromNetwork {
		var acct types.AccountData
		if err := b.c.Get(ctx, "/accounts/"+sender.Hex(), nil, &acct); err != nil {
			return nil, fmt.Errorf("builder: fetch account: %w", err)
		}
		var err error
		seq, err = strconv.ParseUint(acct.SequenceNumber, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("builder: parse sequence_number: %w", err)
		}
	}

	var gasEst types.GasEstimation
	if err := b.c.Get(ctx, "/estimate_gas_price", nil, &gasEst); err != nil {
		return nil, fmt.Errorf("builder: fetch gas price: %w", err)
	}

	maxGas := client.DefaultMaxGasAmount
	gasPrice := gasEst.GasEstimate
	expSecs := uint64(time.Now().Unix()) + client.DefaultTxnExpSecs

	if opts.Options != nil {
		if opts.Options.MaxGasAmount != nil {
			maxGas = *opts.Options.MaxGasAmount
		}
		if opts.Options.GasUnitPrice != nil {
			gasPrice = *opts.Options.GasUnitPrice
		}
		if opts.Options.ExpirationSecs != nil {
			expSecs = uint64(time.Now().Unix()) + *opts.Options.ExpirationSecs
		}
	}

	ef, err := parseEntryFunction(opts.Function, opts.TypeArgs, opts.Args)
	if err != nil {
		return nil, err
	}

	return &RawTransaction{
		Sender:                  sender,
		SequenceNumber:          seq,
		Payload:                 ef,
		MaxGasAmount:            maxGas,
		GasUnitPrice:            gasPrice,
		ExpirationTimestampSecs: expSecs,
		ChainID:                 client.ChainID(b.c.Network()),
	}, nil
}

func parseEntryFunction(function string, typeArgs []string, args [][]byte) (*EntryFunction, error) {
	parts := strings.SplitN(function, "::", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("builder: invalid function format %q, want 0xADDR::module::function", function)
	}
	addr, err := account.AccountAddressFromHex(parts[0])
	if err != nil {
		return nil, fmt.Errorf("builder: invalid module address %q: %w", parts[0], err)
	}

	tags := make([]TypeTag, len(typeArgs))
	for i, ta := range typeArgs {
		tags[i] = ParseTypeTag(ta)
	}

	return &EntryFunction{
		Module:   ModuleID{Address: addr, Name: parts[1]},
		Function: parts[2],
		TypeArgs: tags,
		Args:     args,
	}, nil
}
