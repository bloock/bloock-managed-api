package schema

import (
	"encoding/json"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Certification holds the schema definition for the Certification entity.
type Message struct {
	ent.Schema
}

func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("message").NotEmpty(),
		field.String("root").Default(""),
		field.Int("anchor_id").NonNegative().Default(0),
		field.JSON("proof", json.RawMessage{}).Optional(),
	}
}

// Edges of the Certification.
func (Message) Edges() []ent.Edge {
	return nil
}

func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("message", "root", "anchor_id").
			Unique(),
	}
}
