package response

type IntegrityResponse struct {
	hash     string
	anchorID int
}

func NewIntegrityResponse(hash string, anchorID int) *IntegrityResponse {
	return &IntegrityResponse{
		hash:     hash,
		anchorID: anchorID,
	}
}

func (c IntegrityResponse) Hash() string {
	return c.hash
}

func (c IntegrityResponse) AnchorID() int {
	return c.anchorID
}
