package schema

import "github.com/facebook/ent"

// ReferenceData holds the schema definition for the ReferenceData entity.
type ReferenceData struct {
	ent.Schema
}

// Fields of the ReferenceData.
func (ReferenceData) Fields() []ent.Field {
	return nil
}

// Edges of the ReferenceData.
func (ReferenceData) Edges() []ent.Edge {
	return nil
}
