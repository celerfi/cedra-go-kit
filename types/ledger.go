package types

type LedgerInfo struct {
	ChainID             uint8  `json:"chain_id"`
	Epoch               string `json:"epoch"`
	LedgerVersion       string `json:"ledger_version"`
	OldestLedgerVersion string `json:"oldest_ledger_version"`
	LedgerTimestamp     string `json:"ledger_timestamp"`
	NodeRole            string `json:"node_role"`
	OldestBlockHeight   string `json:"oldest_block_height"`
	BlockHeight         string `json:"block_height"`
	GitHash             string `json:"git_hash"`
}

type GasEstimation struct {
	DeprioritizedGasEstimate *uint64 `json:"deprioritized_gas_estimate,omitempty"`
	GasEstimate              uint64  `json:"gas_estimate"`
	PrioritizedGasEstimate   *uint64 `json:"prioritized_gas_estimate,omitempty"`
}

type Block struct {
	BlockHeight    string        `json:"block_height"`
	BlockHash      string        `json:"block_hash"`
	BlockTimestamp string        `json:"block_timestamp"`
	FirstVersion   string        `json:"first_version"`
	LastVersion    string        `json:"last_version"`
	Transactions   []interface{} `json:"transactions,omitempty"`
}
