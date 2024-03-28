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
		{Name: "proof", Type: field.TypeJSON, Nullable: true},
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
	// MessagesColumns holds the columns for the "messages" table.
	MessagesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "message", Type: field.TypeString},
		{Name: "root", Type: field.TypeString, Default: ""},
		{Name: "anchor_id", Type: field.TypeInt, Default: 0},
		{Name: "proof", Type: field.TypeJSON, Nullable: true},
	}
	// MessagesTable holds the schema information for the "messages" table.
	MessagesTable = &schema.Table{
		Name:       "messages",
		Columns:    MessagesColumns,
		PrimaryKey: []*schema.Column{MessagesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "message_message_root_anchor_id",
				Unique:  true,
				Columns: []*schema.Column{MessagesColumns[1], MessagesColumns[2], MessagesColumns[3]},
			},
		},
	}
	// ProcessesColumns holds the columns for the "processes" table.
	ProcessesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "filename", Type: field.TypeString},
		{Name: "status", Type: field.TypeBool, Default: false},
		{Name: "hash", Type: field.TypeString},
		{Name: "process_response", Type: field.TypeJSON, Nullable: true},
		{Name: "anchor_id", Type: field.TypeInt, Nullable: true},
		{Name: "is_aggregated", Type: field.TypeBool, Default: false},
		{Name: "created_at", Type: field.TypeTime},
	}
	// ProcessesTable holds the schema information for the "processes" table.
	ProcessesTable = &schema.Table{
		Name:       "processes",
		Columns:    ProcessesColumns,
		PrimaryKey: []*schema.Column{ProcessesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "process_id",
				Unique:  true,
				Columns: []*schema.Column{ProcessesColumns[0]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CertificationsTable,
		LocalKeysTable,
		MessagesTable,
		ProcessesTable,
	}
)

func init() {
}
