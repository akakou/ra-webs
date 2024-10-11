package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type TAServer struct {
	ent.Schema
}

// Fields of the TA.
func (TAServer) Fields() []ent.Field {
	return []ent.Field{
		field.String("domain"),
		field.Bytes("public_key"),
		field.String("quote"),
		field.Int("monitor_log_id"),
		field.Bool("is_active").Default(false),
	}
}

// Edges of the TA.
func (TAServer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("violation", TAViolation.Type).Ref("server"),
		edge.To("code", TACode.Type).Unique(),
		edge.To("service", Service.Type).Unique(),
		edge.From("subscription", Subscription.Type).Ref("server"),
	}
}
