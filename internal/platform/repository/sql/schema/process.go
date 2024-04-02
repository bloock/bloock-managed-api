package schema

import (
	"encoding/json"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type Process struct {
	ent.Schema
}

func (Process) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("filename").NotEmpty(),
		field.Bool("status").Default(false),
		field.String("hash").NotEmpty(),
		field.JSON("process_response", json.RawMessage{}).Optional(),
		field.Int("anchor_id").Optional(),
		field.Bool("is_aggregated").Default(false),
		field.Time("created_at").Default(time.Now()),
	}
}

func (Process) Edges() []ent.Edge {
	return nil
}

func (Process) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}
