package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type Subscription struct {
	ent.Schema
}

// Fields of the TA.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.String("subscription"),
	}
}

// Edges of the TA.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("server", TAServer.Type),
	}
}
