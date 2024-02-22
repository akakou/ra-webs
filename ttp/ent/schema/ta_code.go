package schema

import (
	"time"

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
		field.Bytes("unique_id"),
		field.Bytes("public_key"),
		field.String("commit_id"),
		field.Bool("activated").Default(false),
		field.Time("activated_at").Default(time.Now()),
	}
}

// Edges of the TA.
func (TACode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ta", TA.Type),
	}
}
