package availability

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
	"errors"
)

type AvailabilityService struct {
	availabilityRepository repository.AvailabilityRepository
}

func NewAvailabilityService(availabilityRepository repository.AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{availabilityRepository: availabilityRepository}
}

var ErrUnssuportedHosting = errors.New("unsupported hosting type")

func (a AvailabilityService) Upload(ctx context.Context, data []byte, hostingType domain.HostingType) (string, error) {
	switch hostingType {
	case domain.HOSTED:
		return a.availabilityRepository.UploadHosted(ctx, data)
	case domain.IPFS:
		return a.availabilityRepository.UploadIpfs(ctx, data)
	case domain.NONE:
		return "", nil
	default:
		return "", ErrUnssuportedHosting
	}
}
