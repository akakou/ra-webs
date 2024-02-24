package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type TAServer struct {
	ent.Schema
}

// Fields of the TA.
func (TAServer) Fields() []ent.Field {
	return []ent.Field{
		field.String("domain").Unique(),
		field.String("ip").Unique(),
		field.String("service_id"),
		field.Bool("activate").Default(false),
	}
}

// Edges of the TA.
func (TAServer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ta", TA.Type).Ref("server").Unique(),
	}
}
