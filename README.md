# cedra-go-kit

[![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev)
[![CI](https://github.com/celerfi/cedra-go-kit/actions/workflows/ci.yml/badge.svg)](https://github.com/celerfi/cedra-go-kit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/celerfi/cedra-go-kit)](https://goreportcard.com/report/github.com/celerfi/cedra-go-kit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

The first community Go SDK for the [Cedra](https://cedra.dev) network — feature-parity with the official TypeScript SDK.

## Installation

```bash
go get github.com/celerfi/cedra-go-kit
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/celerfi/cedra-go-kit/account"
    "github.com/celerfi/cedra-go-kit/api"
    "github.com/celerfi/cedra-go-kit/client"
    "github.com/celerfi/cedra-go-kit/transaction"
)

func main() {
    ctx := context.Background()

    cedra := api.NewCedra(client.Testnet)

    // Generate a new account
    alice, err := account.GenerateEd25519Account()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Address:", alice.Address().Hex())

    // Fund via faucet (testnet/devnet only)
    if err := cedra.Faucet.FundAccount(ctx, alice.Address().Hex(), 100_000_000); err != nil {
        log.Fatal(err)
    }

    // Check balance
    balance, err := cedra.Account.GetAccountCEDRABalance(ctx, alice.Address().Hex())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %d octas\n", balance)

    // Send CEDRA to another address
    bob, _ := account.GenerateEd25519Account()
    opts := cedra.Coin.TransferCEDRA(bob.Address(), 1_000_000, nil)

    committed, err := cedra.SignAndSubmitTransaction(ctx, alice, opts)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Tx hash:", committed.Hash)
    fmt.Println("Success:", committed.Success)
}
```

## Networks

| Network  | Node API                              | Chain ID |
|----------|---------------------------------------|----------|
| Mainnet  | https://api.mainnet.cedralabs.com/v1  | 1        |
| Testnet  | https://testnet.cedra.dev/v1          | 2        |
| Devnet   | https://devnet.cedra.dev/v1           | 3        |
| Local    | http://127.0.0.1:8080/v1              | 4        |

```go
// Use a specific network
cedra := api.NewCedra(client.Mainnet)

// Or bring your own config
cfg := client.Config{
    Network:    client.Testnet,
    NodeURL:    "https://my-node.example.com/v1",
    IndexerURL: "https://my-indexer.example.com/v1/graphql",
    Timeout:    15 * time.Second,
}
cedra := api.NewCedraWithConfig(cfg)
```

## Account Types

```go
// Ed25519 (default, recommended)
alice, _ := account.GenerateEd25519Account()
alice, _ := account.NewEd25519AccountFromHex("0xYOUR_PRIVATE_KEY")

// SingleKey / Secp256k1
bob, _ := account.GenerateSingleKeyAccount()
bob, _ := account.NewSingleKeyAccountFromHex("0xYOUR_PRIVATE_KEY")
```

## API Reference

| Package | Description |
|---------|-------------|
| `api.GeneralAPI` | Ledger info, gas estimation, block queries, view functions |
| `api.AccountAPI` | Account data, resources, modules, balances, transactions |
| `api.TransactionAPI` | Build, simulate, submit, and wait for transactions |
| `api.EventAPI` | Query events by type, account, or creation number |
| `api.CoinAPI` | CEDRA and custom coin transfer builders |
| `api.FaucetAPI` | Fund accounts on testnet/devnet |
| `api.ANSAPI` | Resolve .cedra names to addresses and vice versa |
| `transaction` | BCS-correct transaction building and signing primitives |
| `bcs` | Binary Canonical Serialization encoder/decoder |
| `crypto` | Ed25519 and Secp256k1 key management |

## Building and Signing Manually

```go
rawTxn, err := cedra.Transaction.BuildTransaction(ctx, alice.Address(), transaction.BuildOptions{
    Function: "0x1::cedra_account::transfer",
    TypeArgs: []string{},
    Args: [][]byte{
        transaction.SerializeAddressArg(bobAddr),
        transaction.SerializeU64Arg(500_000),
    },
})

// Simulate first
simResults, _ := cedra.Transaction.SimulateTransaction(ctx, rawTxn, alice)
fmt.Println("Gas used:", simResults[0].GasUsed)

// Sign and submit
signedBytes, _ := transaction.SignTransaction(rawTxn, alice)
pending, _ := cedra.Transaction.SubmitTransaction(ctx, signedBytes)
committed, _ := cedra.WaitForTransaction(ctx, pending.Hash)
```

## ANS

```go
// Name → address
addr, _ := cedra.ANS.GetAddressFromName(ctx, "alice.cedra")

// Address → name
name, _ := cedra.ANS.GetNameFromAddress(ctx, "0x123...")
```

## License

[MIT](LICENSE) — Celerfi 2026
