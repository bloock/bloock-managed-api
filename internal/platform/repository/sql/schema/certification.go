package schema

import (
	"encoding/json"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Certification holds the schema definition for the Certification entity.
type Certification struct {
	ent.Schema
}

func (Certification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Int("anchor_id").NonNegative(),
		field.String("hash").NotEmpty(),
		field.String("data_id"),
		field.JSON("proof", json.RawMessage{}).Optional(),
	}
}

// Edges of the Certification.
func (Certification) Edges() []ent.Edge {
	return nil
}

func (Certification) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "hash", "anchor_id").
			Unique(),
	}
}
