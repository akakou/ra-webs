package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type TA struct {
	ent.Schema
}

// Fields of the TA.
func (TA) Fields() []ent.Field {
	return []ent.Field{
		field.Bytes("public_key").Nillable(),
		field.Bool("is_active").Default(true),
	}
}

// Edges of the TA.
func (TA) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ct_log", CTLog.Type).Ref("ta"),
		edge.From("evidence_log", EvidenceLog.Type).Ref("ta").Unique(),
	}
}
