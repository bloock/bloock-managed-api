// Code generated by ent, DO NOT EDIT.

package message

import (
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the message type in the database.
	Label = "message"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldMessage holds the string denoting the message field in the database.
	FieldMessage = "message"
	// FieldRoot holds the string denoting the root field in the database.
	FieldRoot = "root"
	// FieldAnchorID holds the string denoting the anchor_id field in the database.
	FieldAnchorID = "anchor_id"
	// FieldProof holds the string denoting the proof field in the database.
	FieldProof = "proof"
	// Table holds the table name of the message in the database.
	Table = "messages"
)

// Columns holds all SQL columns for message fields.
var Columns = []string{
	FieldID,
	FieldMessage,
	FieldRoot,
	FieldAnchorID,
	FieldProof,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// MessageValidator is a validator for the "message" field. It is called by the builders before save.
	MessageValidator func(string) error
	// DefaultRoot holds the default value on creation for the "root" field.
	DefaultRoot string
	// DefaultAnchorID holds the default value on creation for the "anchor_id" field.
	DefaultAnchorID int
	// AnchorIDValidator is a validator for the "anchor_id" field. It is called by the builders before save.
	AnchorIDValidator func(int) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Message queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByMessage orders the results by the message field.
func ByMessage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMessage, opts...).ToFunc()
}

// ByRoot orders the results by the root field.
func ByRoot(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRoot, opts...).ToFunc()
}

// ByAnchorID orders the results by the anchor_id field.
func ByAnchorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAnchorID, opts...).ToFunc()
}
