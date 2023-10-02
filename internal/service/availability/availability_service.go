package availability

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
	"errors"

	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

var ErrUnsupportedHosting = errors.New("unsupported hosting type")

type AvailabilityService struct {
	availabilityRepository repository.AvailabilityRepository
}

func NewAvailabilityService(availabilityRepository repository.AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{availabilityRepository: availabilityRepository}
}

func (a AvailabilityService) Upload(ctx context.Context, record *record.Record, hostingType domain.HostingType) (string, error) {
	switch hostingType {
	case domain.HOSTED:
		hostedID, err := a.availabilityRepository.UploadHosted(ctx, record)
		if err != nil {
			return "", err
		}
		return hostedID, err
	case domain.IPFS:
		ipfsID, err := a.availabilityRepository.UploadIpfs(ctx, record)
		if err != nil {
			return "", err
		}
		return ipfsID, err
	case domain.NONE:
		return "", nil
	default:
		return "", ErrUnsupportedHosting
	}
}

func (a AvailabilityService) Download(ctx context.Context, url string) ([]byte, error) {
	return a.availabilityRepository.FindFile(ctx, url)
}
