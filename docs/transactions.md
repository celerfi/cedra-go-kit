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

## Fee-payer transactions

Use `BuildFeePayerTransaction` for sponsored transactions where a second account pays gas. The signing message follows the fee-payer domain:

```
sha3_256("CEDRA::RawTransactionWithData") || bcs(feePayerRawTransaction)
```

```go
feePayer, _ := account.NewEd25519AccountFromHex("0xFEE_PAYER_PRIVATE_KEY")

feePayerTxn, err := c.Transaction.BuildFeePayerTransaction(ctx, acct.Address(), transaction.BuildOptions{
    Function: "0x1::cedra_account::transfer",
    Args: [][]byte{
        transaction.SerializeAddressArg(recipientAddr),
        transaction.SerializeU64Arg(1_000_000),
    },
    WithFeePayer: true,
})

senderAuthenticator, err := transaction.SignFeePayerTransactionSenderAuthenticator(feePayerTxn, acct)
feePayerAuthenticator, err := transaction.SignFeePayerTransactionFeePayerAuthenticator(feePayerTxn, feePayer)

signedBytes, err := transaction.AssembleFeePayerSignedTransaction(
    feePayerTxn,
    senderAuthenticator,
    feePayer.Address(),
    feePayerAuthenticator,
)
pending, err := c.Transaction.SubmitTransaction(ctx, signedBytes)

// Or simulate with zeroed signatures first.
results, err := c.Transaction.SimulateFeePayerTransaction(ctx, feePayerTxn, acct, feePayer)
fmt.Println(results[0].GasUsed)
```

Browser/backend split:

```go
// Frontend:
feePayerTxn, _ := c.Transaction.BuildFeePayerTransaction(ctx, acct.Address(), transaction.BuildOptions{
    Function:     "0x1::cedra_account::transfer",
    Args:         [][]byte{transaction.SerializeAddressArg(recipientAddr), transaction.SerializeU64Arg(1_000_000)},
    WithFeePayer: true,
})
senderAuthenticator, _ := transaction.SignFeePayerTransactionSenderAuthenticator(feePayerTxn, acct)

// Backend:
feePayerAuthenticator, _ := transaction.SignFeePayerTransactionFeePayerAuthenticator(feePayerTxn, sponsor)
signedBytes, _ := transaction.AssembleFeePayerSignedTransaction(
    feePayerTxn,
    senderAuthenticator,
    sponsor.Address(),
    feePayerAuthenticator,
)
pending, _ := c.Transaction.SubmitTransaction(ctx, signedBytes)
```

## BuildOptions

| Field | Type | Description |
|---|---|---|
| `Function` | `string` | Module function e.g. `0x1::aptos_account::transfer` |
| `TypeArguments` | `[]string` | Generic type params |
| `Arguments` | `[]any` | Function arguments |
| `WithFeePayer` | `bool` | Build a fee-payer/sponsored transaction wrapper |
| `FeePayerAddress` | `*account.AccountAddress` | Optional fee payer address included in the signing wrapper |
| `Options` | `*TransactionOptions` | Optional overrides (gas, expiry, sequence number) |

## TransactionOptions

| Field | Type | Description |
|---|---|---|
| `MaxGasAmount` | `*uint64` | Override default max gas (200,000) |
| `GasUnitPrice` | `*uint64` | Override gas unit price |
| `ExpirationSecs` | `*uint64` | TTL from now in seconds (default 20s) |
| `SequenceNumber` | `*uint64` | **Skip the network fetch and use this value directly.** Use this for high-frequency or concurrent submission — see below. |

## High-frequency / concurrent transactions

By default `BuildTransaction` fetches the account's sequence number from the node on every call. When multiple goroutines call it simultaneously they all receive the **same** sequence number, causing `Transaction already in mempool` collisions.

Manage sequence state in memory and pass it in explicitly to avoid this:

```go
var (
    mu  sync.Mutex
    seq uint64
)

// Fetch once at startup
acctData, _ := cedra.Account.GetAccountInfo(ctx, alice.Address().Hex())
seq, _ = strconv.ParseUint(acctData.SequenceNumber, 10, 64)

func submitNext(payload transaction.BuildOptions) {
    mu.Lock()
    s := seq
    seq++
    mu.Unlock()

    payload.Options = &types.TransactionOptions{SequenceNumber: &s}
    rawTxn, _ := cedra.Transaction.BuildTransaction(ctx, alice.Address(), payload)
    signed, _  := transaction.SignTransaction(rawTxn, alice)
    pending, _ := cedra.Transaction.SubmitTransaction(ctx, signed)
    cedra.WaitForTransaction(ctx, pending.Hash)
}
```

This pattern is the standard approach for keeper bots, activity generators, and any service submitting multiple transactions per second from the same account.

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
