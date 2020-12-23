package gremlin_rest

type ID struct {
	Type string `json:"@type"`
	Value interface{} `json:"@value"`
}
