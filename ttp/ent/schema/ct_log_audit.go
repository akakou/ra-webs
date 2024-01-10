package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CTLog holds the schema definition for the CTLog entity.
type CTLogAudit struct {
	ent.Schema
}

// Fields of the CTLog.
func (CTLogAudit) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_valid"),
		field.String("latest_ct_id"),
	}
}

// Edges of the CTLog.
func (CTLogAudit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ta_info", TAInfo.Type).Unique(),
	}
}
