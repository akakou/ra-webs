package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CTLog holds the schema definition for the CTLog entity.
type CTLog struct {
	ent.Schema
}

// Fields of the CTLog.
func (CTLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("domain"),
		field.Bytes("public_key"),
	}
}

// Edges of the CTLog.
func (CTLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ta_info", TAInfo.Type).Unique(),
	}
}
