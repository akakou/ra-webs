package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type TACode struct {
	ent.Schema
}

// Fields of the TA.
func (TACode) Fields() []ent.Field {
	return []ent.Field{
		field.String("repository"),
		field.String("commit_id"),
		field.Bytes("unique_id"),
		field.Bool("is_active").Default(false),
	}
}

// Edges of the TA.
func (TACode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("server", TAServer.Type).Ref("code"),
		edge.To("service", Service.Type).Unique(),
	}
}
