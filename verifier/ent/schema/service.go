package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type Service struct {
	ent.Schema
}

// Fields of the TA.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("token"),
		field.Bool("is_active").Default(false),
	}
}

// Edges of the TA.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("taserver", TAServer.Type).Ref("service"),
		edge.From("tacode", TACode.Type).Ref("service"),
	}
}
