package domain

import (
	pkg "github.com/bloock/bloock-managed-api/pkg/response"
	"time"
)

type Process struct {
	ID        string    `json:"id"`
	Status    bool      `json:"status"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	pkg.ProcessResponse
}
