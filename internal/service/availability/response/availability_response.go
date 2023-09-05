package response

type AvailabilityResponse struct {
	cid string
}

func NewAvailabilityResponse(cid string) *AvailabilityResponse {
	return &AvailabilityResponse{cid: cid}
}

func (s AvailabilityResponse) Cid() string {
	return s.cid
}
