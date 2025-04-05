package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TAViolation holds the schema definition for the TAViolation entity.
type TAViolation struct {
	ent.Schema
}

// Fields of the TA.
func (TAViolation) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Immutable().Default(func() time.Time {
			return time.Now()
		}),
	}
}

// Edges of the TA.
func (TAViolation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ct_log", CTLog.Type).Unique(),
	}
}
