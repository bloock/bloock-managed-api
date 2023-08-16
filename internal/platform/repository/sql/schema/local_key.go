package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

// LocalKey holds the schema definition for the LocalKey entity.
type LocalKey struct {
	ent.Schema
}
type name struct {
	id       uuid.UUID
	localKey key.LocalKey
	keyType  key.KeyType
}

func (LocalKey) Fields() []ent.Field {

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.JSON("local_key", &key.LocalKey{}),
		field.String("key_type"),
	}
}

// Edges of the LocalKey.
func (LocalKey) Edges() []ent.Edge {
	return nil
}

func (LocalKey) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}
