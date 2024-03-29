// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CertificationsColumns holds the columns for the "certifications" table.
	CertificationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "anchor_id", Type: field.TypeInt},
		{Name: "hash", Type: field.TypeString},
		{Name: "data_id", Type: field.TypeString},
	}
	// CertificationsTable holds the schema information for the "certifications" table.
	CertificationsTable = &schema.Table{
		Name:       "certifications",
		Columns:    CertificationsColumns,
		PrimaryKey: []*schema.Column{CertificationsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "certification_id_hash_anchor_id",
				Unique:  true,
				Columns: []*schema.Column{CertificationsColumns[0], CertificationsColumns[2], CertificationsColumns[1]},
			},
		},
	}
	// LocalKeysColumns holds the columns for the "local_keys" table.
	LocalKeysColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "local_key", Type: field.TypeJSON},
		{Name: "key_type", Type: field.TypeString},
	}
	// LocalKeysTable holds the schema information for the "local_keys" table.
	LocalKeysTable = &schema.Table{
		Name:       "local_keys",
		Columns:    LocalKeysColumns,
		PrimaryKey: []*schema.Column{LocalKeysColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "localkey_id",
				Unique:  true,
				Columns: []*schema.Column{LocalKeysColumns[0]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CertificationsTable,
		LocalKeysTable,
	}
)

func init() {
}
