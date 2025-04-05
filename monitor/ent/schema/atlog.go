package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type ATLog struct {
	ent.Schema
}

// Fields of the TA.
func (ATLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("evidence"),
		field.String("repository"),
		field.String("commit_id"),
		field.Bytes("unique_id"),
		field.Bool("is_active").Default(false),
	}
}

// Edges of the TA.
func (ATLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ct_log", CTLog.Type).Ref("at_log").Unique(),
	}
}
