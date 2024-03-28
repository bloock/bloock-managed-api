// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/certification"
	"github.com/google/uuid"
)

// Certification is the model entity for the Certification schema.
type Certification struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// AnchorID holds the value of the "anchor_id" field.
	AnchorID int `json:"anchor_id,omitempty"`
	// Hash holds the value of the "hash" field.
	Hash string `json:"hash,omitempty"`
	// DataID holds the value of the "data_id" field.
	DataID string `json:"data_id,omitempty"`
	// Proof holds the value of the "proof" field.
	Proof        json.RawMessage `json:"proof,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Certification) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case certification.FieldProof:
			values[i] = new([]byte)
		case certification.FieldAnchorID:
			values[i] = new(sql.NullInt64)
		case certification.FieldHash, certification.FieldDataID:
			values[i] = new(sql.NullString)
		case certification.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Certification fields.
func (c *Certification) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case certification.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				c.ID = *value
			}
		case certification.FieldAnchorID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field anchor_id", values[i])
			} else if value.Valid {
				c.AnchorID = int(value.Int64)
			}
		case certification.FieldHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field hash", values[i])
			} else if value.Valid {
				c.Hash = value.String
			}
		case certification.FieldDataID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field data_id", values[i])
			} else if value.Valid {
				c.DataID = value.String
			}
		case certification.FieldProof:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field proof", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.Proof); err != nil {
					return fmt.Errorf("unmarshal field proof: %w", err)
				}
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Certification.
// This includes values selected through modifiers, order, etc.
func (c *Certification) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// Update returns a builder for updating this Certification.
// Note that you need to call Certification.Unwrap() before calling this method if this Certification
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Certification) Update() *CertificationUpdateOne {
	return NewCertificationClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Certification entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Certification) Unwrap() *Certification {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Certification is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Certification) String() string {
	var builder strings.Builder
	builder.WriteString("Certification(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("anchor_id=")
	builder.WriteString(fmt.Sprintf("%v", c.AnchorID))
	builder.WriteString(", ")
	builder.WriteString("hash=")
	builder.WriteString(c.Hash)
	builder.WriteString(", ")
	builder.WriteString("data_id=")
	builder.WriteString(c.DataID)
	builder.WriteString(", ")
	builder.WriteString("proof=")
	builder.WriteString(fmt.Sprintf("%v", c.Proof))
	builder.WriteByte(')')
	return builder.String()
}

// Certifications is a parsable slice of Certification.
type Certifications []*Certification
