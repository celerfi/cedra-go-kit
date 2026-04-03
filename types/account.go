package types

import "encoding/json"

type AccountData struct {
	SequenceNumber    string `json:"sequence_number"`
	AuthenticationKey string `json:"authentication_key"`
}

type MoveResource struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MoveModuleBytecode struct {
	Bytecode string      `json:"bytecode"`
	Abi      interface{} `json:"abi,omitempty"`
}

type CoinStore struct {
	Coin struct {
		Value string `json:"value"`
	} `json:"coin"`
}
