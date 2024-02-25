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
		field.Bytes("public_key"),
		field.Bool("is_valid").Default(false),
	}
}

// Edges of the TA.
func (TA) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("code", TACode.Type).Unique(),
		edge.To("server", TAServer.Type).Unique(),
		edge.To("ct_audit", CTAudit.Type).Unique(),
	}
}
