package response

import (
	"fmt"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
)

type AvailabilityResponse struct {
	ty          string
	id          string
	url         string
	contentType string
	size        int64
}

func NewAvailabilityResponse(id string, hostingType domain.HostingType, contentType string, size int64) *AvailabilityResponse {
	var url string
	switch hostingType {
	case domain.IPFS:
		url = fmt.Sprintf("%s/hosting/v1/ipfs/%s", config.Configuration.Bloock.CdnHost, id)
	case domain.HOSTED:
		url = fmt.Sprintf("%s/hosting/v1/hosted/%s", config.Configuration.Bloock.CdnHost, id)
	default:
		url = ""
	}
	return &AvailabilityResponse{
		ty:          hostingType.String(),
		id:          id,
		url:         url,
		contentType: contentType,
		size:        size,
	}
}

func (s AvailabilityResponse) Type() string {
	return s.ty
}

func (s AvailabilityResponse) Id() string {
	return s.id
}

func (s AvailabilityResponse) Url() string {
	return s.url
}

func (s AvailabilityResponse) ContentType() string {
	return s.contentType
}

func (s AvailabilityResponse) Size() int64 {
	return s.size
}
