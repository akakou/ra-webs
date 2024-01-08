package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type TA struct {
	ent.Schema
}

// Fields of the TA.
func (TA) Fields() []ent.Field {
	return []ent.Field{
		field.String("domain"),
		field.String("public_key"),
		field.String("attestation"),
	}
}

// Edges of the TA.
func (TA) Edges() []ent.Edge {
	return nil
}
