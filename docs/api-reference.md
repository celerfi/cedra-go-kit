---
title: API Reference
nav_order: 4
---

# API Reference

All APIs are accessed through the top-level client:

```go
c := cedra.New(client.Testnet)

c.Account     // AccountAPI
c.Transaction // TransactionAPI
c.Faucet      // FaucetAPI
c.General     // GeneralAPI
c.Event       // EventAPI
c.Coin        // CoinAPI
```

---

## AccountAPI

| Method | Description |
|---|---|
| `GetAccountInfo(ctx, address)` | Sequence number and auth key |
| `GetSequenceNumber(ctx, address)` | Current sequence number as `uint64` |
| `GetAccountResources(ctx, address)` | All Move resources on an account |
| `GetAccountResource(ctx, address, type)` | Single resource by type tag |
| `GetAccountModules(ctx, address)` | All deployed Move modules |
| `GetAccountTransactions(ctx, address, limit, start)` | Paginated transaction history |
| `GetAccountCEDRABalance(ctx, address)` | CEDRA balance via FA standard |
| `GetAccountFABalance(ctx, address, faAddress)` | Any fungible asset balance |

---

## TransactionAPI

| Method | Description |
|---|---|
| `BuildTransaction(ctx, sender, opts)` | Returns a `RawTransaction` ready to sign |
| `SimulateTransaction(ctx, rawTxn, signer)` | Dry-run with zeroed signature |
| `SubmitTransaction(ctx, signedBCS)` | Submit signed BCS bytes |
| `WaitForTransaction(ctx, hash)` | Poll until committed or timeout |
| `GetTransactionByHash(ctx, hash)` | Fetch committed transaction |
| `GetTransactionByVersion(ctx, version)` | Fetch by ledger version |
| `GetTransactions(ctx, limit, start)` | Paginated recent transactions |
| `IsTransactionPending(ctx, hash)` | `true` if not yet committed |

---

## FaucetAPI

| Method | Description |
|---|---|
| `FundAccount(ctx, address, amount)` | Fund and wait for confirmation |
| `FundAccountNoWait(ctx, address, amount)` | Fund and return tx hashes immediately |

Only available on testnet and devnet.

---

## GeneralAPI

| Method | Description |
|---|---|
| `GetLedgerInfo(ctx)` | Chain ID, ledger version, timestamp |
| `GetChainID(ctx)` | Chain ID as `uint8` |
| `HealthCheck(ctx)` | Node liveness check |
| `ViewFunction(ctx, req)` | Call any Move view function |

---

## EventAPI

| Method | Description |
|---|---|
| `GetEventsByEventHandle(ctx, address, handle, field, limit, start)` | Paginated events by handle |
| `GetEventsByCreationNumber(ctx, address, creationNumber, limit, start)` | Paginated events by creation number |

---

## CoinAPI

| Method | Description |
|---|---|
| `Transfer(ctx, sender, recipient, amount)` | `0x1::aptos_account::transfer` shorthand |

---

## Custom network config

```go
import "github.com/celerfi/cedra-go-kit/client"

cfg := client.Config{
    NodeURL:   "https://my-node.example.com/v1",
    FaucetURL: "https://my-faucet.example.com",
    Timeout:   15 * time.Second,
}
c := cedra.NewWithConfig(cfg)
```

---

## Error handling

All methods return `error`. API errors carry the HTTP status code and response body:

```go
txn, err := c.Transaction.GetTransactionByHash(ctx, hash)
if err != nil {
    var apiErr *client.APIError
    if errors.As(err, &apiErr) {
        fmt.Println(apiErr.StatusCode) // 404
        fmt.Println(apiErr.Message)
    }
}
```
