---
title: Transactions
nav_order: 3
---

# Transactions

## Build, sign, and submit

```go
import (
    cedra "github.com/celerfi/cedra-go-kit"
    "github.com/celerfi/cedra-go-kit/account"
    "github.com/celerfi/cedra-go-kit/client"
    "github.com/celerfi/cedra-go-kit/transaction"
)

c := cedra.New(client.Testnet)
acct, _ := account.NewEd25519Account()

// 1. Build
rawTxn, err := c.Transaction.BuildTransaction(ctx, acct.Address(), transaction.BuildOptions{
    Function:      "0x1::aptos_account::transfer",
    TypeArguments: []string{},
    Arguments:     []any{recipientAddr.Hex(), "1000000"},
})

// 2. Sign
signedBytes, err := transaction.SignTransaction(rawTxn, acct)

// 3. Submit
pending, err := c.Transaction.SubmitTransaction(ctx, signedBytes)

// 4. Wait for confirmation
committed, err := c.Transaction.WaitForTransaction(ctx, pending.Hash)
fmt.Println(committed.Success) // true
```

## Simulate before submitting

Simulation lets you check gas cost and output without spending tokens.

```go
results, err := c.Transaction.SimulateTransaction(ctx, rawTxn, acct)
fmt.Println(results[0].GasUsed)
```

The simulation endpoint requires a zeroed signature (64 zero bytes). The SDK handles this automatically via `transaction.SimulateTransaction`.

## BuildOptions

| Field | Type | Description |
|---|---|---|
| `Function` | `string` | Module function e.g. `0x1::aptos_account::transfer` |
| `TypeArguments` | `[]string` | Generic type params |
| `Arguments` | `[]any` | Function arguments |
| `MaxGasAmount` | `uint64` | Override default (200,000) |
| `GasUnitPrice` | `uint64` | Override gas price |
| `ExpirationSecs` | `uint64` | TTL from now (default 20s) |

## Signing internals

The signing message follows the BCS standard:

```
sha3_256("CEDRA::RawTransaction") || bcs(rawTransaction)
```

Public keys and signatures are **length-prefixed** (ULEB128) in the authenticator, not fixed bytes. The SDK handles this correctly for both Ed25519 and Secp256k1.

## Fetch by hash or version

```go
txn, err := c.Transaction.GetTransactionByHash(ctx, "0xabc...")
txn, err := c.Transaction.GetTransactionByVersion(ctx, 12345)
```

## Check pending status

```go
pending, err := c.Transaction.IsTransactionPending(ctx, hash)
```
