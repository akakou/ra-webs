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
		field.Bytes("public_key"),
		field.Int("monitor_log_id"),
		field.Bool("is_active").Default(false),
	}
}

// Edges of the TA.
func (CTLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("violation", Violation.Type).Ref("ct_log"),
		edge.To("at_log", ATLog.Type).Unique(),
		edge.From("subscription", Subscription.Type).Ref("ct_log"),
	}
}
