// Code generated by ent, DO NOT EDIT.

package ent

import (
	"bloock-managed-api/internal/platform/repository/sql/ent/localkey"
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

// LocalKey is the model entity for the LocalKey schema.
type LocalKey struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// LocalKey holds the value of the "local_key" field.
	LocalKey *key.LocalKey `json:"local_key,omitempty"`
	// KeyType holds the value of the "key_type" field.
	KeyType      string `json:"key_type,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LocalKey) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case localkey.FieldLocalKey:
			values[i] = new([]byte)
		case localkey.FieldKeyType:
			values[i] = new(sql.NullString)
		case localkey.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LocalKey fields.
func (lk *LocalKey) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case localkey.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				lk.ID = *value
			}
		case localkey.FieldLocalKey:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field local_key", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &lk.LocalKey); err != nil {
					return fmt.Errorf("unmarshal field local_key: %w", err)
				}
			}
		case localkey.FieldKeyType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field key_type", values[i])
			} else if value.Valid {
				lk.KeyType = value.String
			}
		default:
			lk.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the LocalKey.
// This includes values selected through modifiers, order, etc.
func (lk *LocalKey) Value(name string) (ent.Value, error) {
	return lk.selectValues.Get(name)
}

// Update returns a builder for updating this LocalKey.
// Note that you need to call LocalKey.Unwrap() before calling this method if this LocalKey
// was returned from a transaction, and the transaction was committed or rolled back.
func (lk *LocalKey) Update() *LocalKeyUpdateOne {
	return NewLocalKeyClient(lk.config).UpdateOne(lk)
}

// Unwrap unwraps the LocalKey entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (lk *LocalKey) Unwrap() *LocalKey {
	_tx, ok := lk.config.driver.(*txDriver)
	if !ok {
		panic("ent: LocalKey is not a transactional entity")
	}
	lk.config.driver = _tx.drv
	return lk
}

// String implements the fmt.Stringer.
func (lk *LocalKey) String() string {
	var builder strings.Builder
	builder.WriteString("LocalKey(")
	builder.WriteString(fmt.Sprintf("id=%v, ", lk.ID))
	builder.WriteString("local_key=")
	builder.WriteString(fmt.Sprintf("%v", lk.LocalKey))
	builder.WriteString(", ")
	builder.WriteString("key_type=")
	builder.WriteString(lk.KeyType)
	builder.WriteByte(')')
	return builder.String()
}

// LocalKeys is a parsable slice of LocalKey.
type LocalKeys []*LocalKey