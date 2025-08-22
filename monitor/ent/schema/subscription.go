package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type Subscription struct {
	ent.Schema
}

// Fields of the TA.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.String("endpoint"),
		field.String("p256dh"),
		field.String("auth"),
	}
}

// Edges of the TA.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{}
}
