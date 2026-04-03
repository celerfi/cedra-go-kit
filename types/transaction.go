package types

type PendingTransaction struct {
	Hash                    string `json:"hash"`
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	Type                    string `json:"type"`
}

type CommittedTransaction struct {
	Hash                    string  `json:"hash"`
	Sender                  string  `json:"sender"`
	SequenceNumber          string  `json:"sequence_number"`
	MaxGasAmount            string  `json:"max_gas_amount"`
	GasUnitPrice            string  `json:"gas_unit_price"`
	ExpirationTimestampSecs string  `json:"expiration_timestamp_secs"`
	GasUsed                 string  `json:"gas_used"`
	Success                 bool    `json:"success"`
	VmStatus                string  `json:"vm_status"`
	Version                 string  `json:"version"`
	Events                  []Event `json:"events"`
	Timestamp               string  `json:"timestamp"`
	Type                    string  `json:"type"`
}

type ViewRequest struct {
	Function      string `json:"function"`
	TypeArguments []string `json:"type_arguments"`
	Arguments     []any    `json:"arguments"`
}

type TransactionOptions struct {
	MaxGasAmount   *uint64
	GasUnitPrice   *uint64
	ExpirationSecs *uint64
}
