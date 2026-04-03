package client

type Network string

const (
	Mainnet Network = "mainnet"
	Testnet Network = "testnet"
	Devnet  Network = "devnet"
	Local   Network = "local"
)

var networkChainIDs = map[Network]uint8{
	Mainnet: 1,
	Testnet: 2,
	Devnet:  3,
	Local:   4,
}

type networkEndpoints struct {
	Node    string
	Indexer string
	Faucet  string
}

var endpoints = map[Network]networkEndpoints{
	Mainnet: {
		Node:    "https://api.mainnet.cedralabs.com/v1",
		Indexer: "https://graphql.cedra.dev/v1/graphql",
		Faucet:  "",
	},
	Testnet: {
		Node:    "https://testnet.cedra.dev/v1",
		Indexer: "https://graphql.cedra.dev/v1/graphql",
		Faucet:  "https://faucet-api.cedra.dev",
	},
	Devnet: {
		Node:    "https://devnet.cedra.dev/v1",
		Indexer: "https://graphql-devnet.cedra.dev/v1/graphql",
		Faucet:  "https://devfaucet-api.cedra.dev",
	},
	Local: {
		Node:    "http://127.0.0.1:8080/v1",
		Indexer: "http://127.0.0.1:8090/v1/graphql",
		Faucet:  "http://127.0.0.1:8081",
	},
}

func ChainID(n Network) uint8 {
	return networkChainIDs[n]
}
