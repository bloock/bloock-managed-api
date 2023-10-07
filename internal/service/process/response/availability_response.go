package response

import (
	"bloock-managed-api/internal/domain"
	"fmt"
)

type AvailabilityResponse struct {
	id  string
	url string
}

func NewAvailabilityResponse(id string, hostingType domain.HostingType) *AvailabilityResponse {
	var url string
	switch hostingType {
	case domain.IPFS:
		url = fmt.Sprintf("https://cdn.bloock.com/hosting/v1/ipfs/%s", id)
	case domain.HOSTED:
		url = fmt.Sprintf("https://cdn.bloock.com/hosting/v1/hosted/%s", id)
	default:
		url = ""
	}
	return &AvailabilityResponse{
		id:  id,
		url: url,
	}
}

func (s AvailabilityResponse) Id() string {
	return s.id
}

func (s AvailabilityResponse) Url() string {
	return s.url
}
