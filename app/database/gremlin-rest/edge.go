package gremlin_rest

type EdgeValue struct {
	ID ID `json:"id"`
	Label string `json:"label"`
	InVLabel string `json:"inVLabel"`
	OutVLabel string `json:"outVLabel"`
	InV ID `json:"inV"`
	OutV ID `json:"outV"`
}
type Edge struct {
	Type  string    `json:"@type"`
	Value EdgeValue `json:"@value"`
}
type EdgeList struct {
	Type string  `json:"@type"`
	Value []Edge `json:"@value"`
}
