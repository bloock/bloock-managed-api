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
		hostedID, err := a.availabilityRepository.UploadHosted(ctx, data)
		if err != nil {
			return "", err
		}
		return hostedID, err
	case domain.IPFS:
		ipfsID, err := a.availabilityRepository.UploadIpfs(ctx, data)
		if err != nil {
			return "", err
		}
		return ipfsID, err
	case domain.NONE:
		return "", nil
	default:
		return "", ErrUnssuportedHosting
	}
}
