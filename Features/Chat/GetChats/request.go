package getchats

type request struct {
	ExclusiveStartKey uint64 `json:"exclusiveStartKey"`
	Limit             uint   `json:"limit"`
}
