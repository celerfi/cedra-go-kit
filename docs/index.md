---
title: Home
nav_order: 1
---

# cedra-go-kit

The community Go SDK for the [Cedra](https://cedra.dev) blockchain. Feature-parity with the official TypeScript SDK.

[![Go](https://img.shields.io/badge/go-1.25-blue)](https://go.dev)
[![CI](https://github.com/celerfi/cedra-go-kit/actions/workflows/ci.yml/badge.svg)](https://github.com/celerfi/cedra-go-kit/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Install

```bash
go get github.com/celerfi/cedra-go-kit
```

Requires Go 1.25+.

---

## Quick start

```go
package main

import (
    "context"
    "fmt"

    cedra "github.com/celerfi/cedra-go-kit"
    "github.com/celerfi/cedra-go-kit/account"
    "github.com/celerfi/cedra-go-kit/client"
)

func main() {
    c := cedra.New(client.Testnet)

    // Generate a new Ed25519 account
    acct, _ := account.NewEd25519Account()
    addr := acct.Address().Hex()

    // Fund from faucet
    c.Faucet.FundAccount(context.Background(), addr, 100_000_000)

    // Check balance
    bal, _ := c.Account.GetAccountCEDRABalance(context.Background(), addr)
    fmt.Printf("balance: %d\n", bal)
}
```

---

## Networks

| Constant | Node URL |
|---|---|
| `client.Testnet` | `https://testnet.cedra.dev/v1` |
| `client.Devnet` | `https://devnet.cedra.dev/v1` |

---

## Packages

| Package | Purpose |
|---|---|
| [`account`](./accounts) | Key generation, address derivation |
| [`api`](./api-reference) | Node REST API calls |
| [`transaction`](./transactions) | Build, sign, simulate, submit |
| `client` | HTTP client, network config |
| `bcs` | Binary Canonical Serialization |
| `crypto` | Ed25519 and Secp256k1 primitives |
