package availability

import (
	"bloock-managed-api/internal/domain/repository"
	hosting "bloock-managed-api/internal/service"
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

func (a AvailabilityService) Upload(ctx context.Context, data []byte, hostingType hosting.HostingType) (string, error) {
	switch hostingType {
	case hosting.HOSTED:
		return a.availabilityRepository.UploadHosted(ctx, data)
	case hosting.IPFS:
		return a.availabilityRepository.UploadIpfs(ctx, data)
	case hosting.NONE:
		return "", nil
	default:
		return "", ErrUnssuportedHosting
	}
}
