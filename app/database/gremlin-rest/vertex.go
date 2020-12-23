package gremlin_rest

type VertexValue struct {
	ID      ID     `json:"id"`
	Label string `json:"label"`
}
type Vertex struct {
	Type  string `json:"@type"`
	Value VertexValue `json:"@value"`
}
type VertexList struct {
	Type string                 `json:"@type"`
	Value     []Vertex `json:"@value"`
}
