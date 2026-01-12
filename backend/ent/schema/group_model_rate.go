package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// GroupModelRate holds the schema definition for per-model rate configuration within a group.
type GroupModelRate struct {
	ent.Schema
}

func (GroupModelRate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "group_model_rates"},
	}
}

func (GroupModelRate) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("group_id"),
		field.String("model").
			MaxLen(100).
			NotEmpty(),
		field.Float("rate_multiplier").
			SchemaType(map[string]string{dialect.Postgres: "decimal(10,4)"}).
			Default(1.0),
		field.Float("card_price").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Optional().
			Nillable().
			Comment("次卡模式单次请求价格(USD)，NULL表示不支持次卡"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (GroupModelRate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).
			Ref("model_rates").
			Unique().
			Required().
			Field("group_id"),
	}
}

func (GroupModelRate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("group_id", "model").Unique(),
		index.Fields("group_id"),
		index.Fields("model"),
	}
}
