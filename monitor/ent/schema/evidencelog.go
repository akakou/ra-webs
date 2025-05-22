package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type EvidenceLog struct {
	ent.Schema
}

// Fields of the TA.
func (EvidenceLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("evidence"),
		field.String("repository"),
		field.String("commit_id"),
		field.Bytes("unique_id"),
	}
}

// Edges of the TA.
func (EvidenceLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ta", TA.Type).Unique(),
	}
}
