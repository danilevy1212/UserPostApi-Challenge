package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("title").
			NotEmpty(),
		field.String("content").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Uint64("user_id").
			Positive().
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("posts").
			Field("user_id").
			Unique().
			Immutable().
			Required(),
	}
}
