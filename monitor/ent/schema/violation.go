package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TAViolation holds the schema definition for the TAViolation entity.
type Violation struct {
	ent.Schema
}

// Fields of the TA.
func (Violation) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Immutable().Default(func() time.Time {
			return time.Now()
		}),
	}
}

// Edges of the TA.
func (Violation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ta", TA.Type).Unique(),
	}
}
