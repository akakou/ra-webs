package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type CTLog struct {
	ent.Schema
}

// Fields of the TA.
func (CTLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("monitor_log_id"),
	}
}

// Edges of the TA.
func (CTLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ta", TA.Type).Unique(),
	}
}
