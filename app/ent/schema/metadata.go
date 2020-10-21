package schema

import "github.com/facebook/ent"

// MetaData holds the schema definition for the MetaData entity.
type MetaData struct {
	ent.Schema
}

// Fields of the MetaData.
func (MetaData) Fields() []ent.Field {
	return nil
}

// Edges of the MetaData.
func (MetaData) Edges() []ent.Edge {
	return nil
}
