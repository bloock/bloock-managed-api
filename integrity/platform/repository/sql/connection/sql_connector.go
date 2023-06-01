package connection

import (
	"bloock-managed-api/ent"
)

type SQLConnector interface {
	Connect(driver string, connectionURL string) (*ent.Client, error)
}
