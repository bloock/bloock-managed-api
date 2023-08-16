package service

import (
	"bloock-managed-api/internal/domain"
	create_request "bloock-managed-api/internal/service/create/request"
	create_response "bloock-managed-api/internal/service/create/response"
	update_request "bloock-managed-api/internal/service/update/request"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
)

type SignService interface {
	Sign(ctx context.Context, req create_request.SignRequest) (*create_response.SignResponse, error)
}

type CertificateService interface {
	Certify(ctx context.Context, files [][]byte) ([]create_response.CertificationResponse, error)
}

type CertificateUpdateAnchorService interface {
	UpdateAnchor(ctx context.Context, updateRequest update_request.UpdateCertificationAnchorRequest) error
}

type LocalKeyCreateService interface {
	Create(ctx context.Context, keyType key.KeyType) (domain.LocalKey, error)
}

type ManagedKeyCreateService interface {
	Create(request create_request.CreateManagedKeyRequest) (key.ManagedKey, error)
}

type GetLocalKeysService interface {
	Get(ctx context.Context) ([]domain.LocalKey, error)
}
