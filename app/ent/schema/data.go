package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Data holds the schema definition for the Data entity.
type Data struct {
	ent.Schema
}

// Fields of the Data.
func (Data) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default("unknown"),
		field.String("context").Default("unknown"),
	}
}

// Edges of the Data.
func (Data) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("has_table", Data.Type),
		edge.To("has_field", Data.Type),
	}
}
