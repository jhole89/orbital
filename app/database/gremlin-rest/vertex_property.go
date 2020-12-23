package gremlin_rest

type PropertyValue struct {
	ID    ID     `json:"@id"`
	Value string `json:"value"`
	Label string `json:"label"`
}
type VertexProperty struct {
	Type  string         `json:"@type"`
	Value PropertyValue `json:"@value"`
}
type VertexPropertyList struct {
	Type string                     `json:"@type"`
	Value []VertexProperty `json:"@value"`
}
