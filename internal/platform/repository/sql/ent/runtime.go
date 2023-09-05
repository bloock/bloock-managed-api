// Code generated by ent, DO NOT EDIT.

package ent

import (
	"bloock-managed-api/internal/platform/repository/sql/ent/certification"
	"bloock-managed-api/internal/platform/repository/sql/ent/localkey"
	"bloock-managed-api/internal/platform/repository/sql/schema"

	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	certificationFields := schema.Certification{}.Fields()
	_ = certificationFields
	// certificationDescAnchorID is the schema descriptor for anchor_id field.
	certificationDescAnchorID := certificationFields[1].Descriptor()
	// certification.AnchorIDValidator is a validator for the "anchor_id" field. It is called by the builders before save.
	certification.AnchorIDValidator = certificationDescAnchorID.Validators[0].(func(int) error)
	// certificationDescHash is the schema descriptor for hash field.
	certificationDescHash := certificationFields[2].Descriptor()
	// certification.HashValidator is a validator for the "hash" field. It is called by the builders before save.
	certification.HashValidator = certificationDescHash.Validators[0].(func(string) error)
	// certificationDescID is the schema descriptor for id field.
	certificationDescID := certificationFields[0].Descriptor()
	// certification.DefaultID holds the default value on creation for the id field.
	certification.DefaultID = certificationDescID.Default.(func() uuid.UUID)
	localkeyFields := schema.LocalKey{}.Fields()
	_ = localkeyFields
	// localkeyDescID is the schema descriptor for id field.
	localkeyDescID := localkeyFields[0].Descriptor()
	// localkey.DefaultID holds the default value on creation for the id field.
	localkey.DefaultID = localkeyDescID.Default.(func() uuid.UUID)
}
