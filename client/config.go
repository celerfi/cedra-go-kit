package client

import "time"

const (
	DefaultMaxGasAmount   uint64 = 200_000
	DefaultTxnExpSecs     uint64 = 20
	DefaultTxnTimeoutSecs uint64 = 20
	CedraCoin                    = "0x1::cedra_coin::CedraCoin"
	CedraFA                      = "0x000000000000000000000000000000000000000000000000000000000000000a"
)

type Config struct {
	Network    Network
	NodeURL    string
	IndexerURL string
	FaucetURL  string
	Timeout    time.Duration
}

func DefaultConfig(network Network) Config {
	ep := endpoints[network]
	return Config{
		Network:    network,
		NodeURL:    ep.Node,
		IndexerURL: ep.Indexer,
		FaucetURL:  ep.Faucet,
		Timeout:    30 * time.Second,
	}
}
