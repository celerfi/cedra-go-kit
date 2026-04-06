---
title: Accounts
nav_order: 2
---

# Accounts

## Generate an account

```go
import "github.com/celerfi/cedra-go-kit/account"

// Ed25519 (default, recommended)
acct, err := account.NewEd25519Account()

// Secp256k1
acct, err := account.NewSingleKeyAccount()
```

Both types implement the `Account` interface:

```go
type Account interface {
    Address() AccountAddress
    AuthKey() []byte
    PublicKeyBytes() []byte
    SignTransaction(signingMessage []byte) ([]byte, error)
}
```

## Load from private key

```go
privKeyHex := "0xabc123..."
acct, err := account.Ed25519AccountFromHex(privKeyHex)
```

## AccountAddress

```go
// From hex string
addr, err := account.AccountAddressFromHex("0x1")

// To hex
fmt.Println(acct.Address().Hex()) // "0x000...001"
```

Addresses are always padded to 32 bytes internally. `AccountAddressFromHex` handles both short (`0x1`) and full-length (`0x000...001`) forms.

## Fund from faucet (testnet/devnet only)

```go
err := c.Faucet.FundAccount(ctx, acct.Address().Hex(), 100_000_000)
```

`FundAccount` waits for the faucet transaction to land on-chain before returning. Use `FundAccountNoWait` if you want the hash back immediately.

## Get CEDRA balance

```go
bal, err := c.Account.GetAccountCEDRABalance(ctx, addr)
```

Cedra uses the [Fungible Asset](https://cedra.dev) standard. Balance is fetched via a view function on `0x1::primary_fungible_store`.

## Get any FA balance

```go
faAddress := "0x000...000a" // CEDRA FA metadata address
bal, err := c.Account.GetAccountFABalance(ctx, addr, faAddress)
```
