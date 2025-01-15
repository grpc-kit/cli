package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Demo holds the schema definition for the Demo entity.
type Demo struct {
	ent.Schema
}

// Fields of the Demo.
func (Demo) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default("grpc-kit"),
	}
}

// Edges of the Demo.
func (Demo) Edges() []ent.Edge {
	return nil
}

// Mixin of the Demo.
func (Demo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Annotations 自定义表名
func (Demo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_demo_can_remove"},
	}
}
