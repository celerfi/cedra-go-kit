package types

import "encoding/json"

type Event struct {
	Version        string          `json:"version"`
	GUID           EventGUID       `json:"guid"`
	SequenceNumber string          `json:"sequence_number"`
	Type           string          `json:"type"`
	Data           json.RawMessage `json:"data"`
}

type EventGUID struct {
	CreationNumber string `json:"creation_number"`
	AccountAddress string `json:"account_address"`
}
