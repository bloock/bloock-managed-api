package create

import (
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/create/request"
	"bloock-managed-api/internal/service/create/response"
	"context"
)

type SignatureService struct {
	authenticityRepository repository.AuthenticityRepository
	localKeysRepository    repository.LocalKeysRepository
}

func NewSignature(authenticityRepository repository.AuthenticityRepository, localKeysRepository repository.LocalKeysRepository) *SignatureService {
	return &SignatureService{authenticityRepository: authenticityRepository, localKeysRepository: localKeysRepository}
}

func (s SignatureService) Sign(ctx context.Context, req request.SignRequest) (*response.SignResponse, error) {

	localKey, err := s.localKeysRepository.FindKeyByID(ctx, req.LocalKey())
	if err != nil {
		return nil, err
	}
	signedRecord, err := s.authenticityRepository.Sign(
		localKey.LocalKey(),
		localKey.KeyType(),
		req.CommonName(),
		req.Data(),
	)
	if err != nil {
		return nil, err
	}

	return response.NewSignResponse(signedRecord.Retrieve()), nil
}
