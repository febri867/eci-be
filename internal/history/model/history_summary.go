package model

type HistorySummary struct {
	Total   int64 `json:"total"`
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
	Deleted int64 `json:"deleted"`
}
