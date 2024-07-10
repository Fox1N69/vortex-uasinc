package models

type AlgorithmStatus struct {
	ID       int64 `json:"id"`
	ClientID int64 `json:"client_id"`
	VWAP     bool  `json:"vwap"`
	TWAP     bool  `json:"twap"`
	HFT      bool  `json:"hft"`
}
