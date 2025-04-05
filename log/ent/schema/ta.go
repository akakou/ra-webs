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
		field.String("evidence"),
		field.Bytes("signature"),
		field.String("repository"),
		field.String("commit_id"),
	}
}

// Edges of the TA.
func (TA) Edges() []ent.Edge {
	return nil
}
