---
title: BCS
nav_order: 5
---

# BCS (Binary Canonical Serialization)

The `bcs` package implements the serialization format used by all Move transactions on Cedra.

You won't need this directly for most use cases — the SDK handles BCS internally when building and signing transactions.

## Serializer

```go
import "github.com/celerfi/cedra-go-kit/bcs"

s := &bcs.Serializer{}
s.SerializeU64(12345)
s.SerializeBytes([]byte{0x01, 0x02})  // length-prefixed (ULEB128 + bytes)
s.SerializeFixedBytes([]byte{0x01})   // no length prefix
s.SerializeULEB128(200_000)

raw := s.ToBytes()
```

## Deserializer

```go
d := bcs.NewDeserializer(raw)
n := d.DeserializeU64()
b := d.DeserializeBytes()
```

## Key serialization rules

- **Public keys and signatures** use `SerializeBytes` (length-prefixed). Using `SerializeFixedBytes` here will cause a deserialize error on the node.
- **Account addresses** use `SerializeFixedBytes` (always 32 bytes, no prefix).
- **Integers** use their fixed-size variants (`SerializeU8`, `SerializeU32`, `SerializeU64`).
