package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TA holds the schema definition for the TA entity.
type TAInfo struct {
	ent.Schema
}

// Fields of the TA.
func (TAInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("domain"),
		field.String("ip_address"),
		field.String("git_repository"),
	}
}

// Edges of the TA.
func (TAInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ct_log", CTLogAudit.Type).Ref("ta_info").Unique(),
		edge.From("ta_code", TACode.Type).Ref("ta_info"),
	}
}
