package availability

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrUnsupportedHosting = errors.New("unsupported hosting type")

type AvailabilityService struct {
	availabilityRepository repository.AvailabilityRepository
}

func NewAvailabilityService(availabilityRepository repository.AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{availabilityRepository: availabilityRepository}
}

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
		return "", ErrUnsupportedHosting
	}
}

func (a AvailabilityService) Download(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("error downloading file from %s: %s", url, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("error downloading file from %s: received status code %d", url, resp.StatusCode)
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error downloading file from %s: %s", url, err.Error())
	}

	return file, nil

}
