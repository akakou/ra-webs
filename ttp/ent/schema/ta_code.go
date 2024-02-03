package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type TACode struct {
	ent.Schema
}

// Fields of the TA.
func (TACode) Fields() []ent.Field {
	return []ent.Field{
		field.String("unique_id"),
		field.String("commit_id"),
	}
}

// Edges of the TA.
func (TACode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ta_info", TAInfo.Type),
	}
}
