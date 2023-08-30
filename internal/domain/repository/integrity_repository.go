package repository

import (
	"bloock-managed-api/internal/domain"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
)

type IntegrityRepository interface {
	Certify(ctx context.Context, bytes []byte) (certification domain.Certification, err error)
	GetAnchorByID(ctx context.Context, anchorID int) (integrity.Anchor, error)
}
