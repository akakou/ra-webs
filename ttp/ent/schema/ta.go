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
		field.String("domain"),
		field.String("ip"),
		field.String("git"),
	}
}

// Edges of the TA.
func (TA) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("audit_log", TAAuditLog.Type).Ref("ta").Unique(),
		edge.From("code", TACode.Type).Ref("ta"),
	}
}
