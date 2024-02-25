package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TACTAudit entity.
type CTAudit struct {
	ent.Schema
}

// Fields of the TACTAudit.
func (CTAudit) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_valid").Default(true),
		field.String("last"),
	}
}

// Edges of the TACTAudit.
func (CTAudit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ta", TA.Type).Ref("ct_audit"),
	}
}
