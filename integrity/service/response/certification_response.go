package response

type CertificationResponse struct {
	hash     string
	anchorID int
}

func NewCertificationResponse(hash string, anchorID int) *CertificationResponse {
	return &CertificationResponse{
		hash:     hash,
		anchorID: anchorID,
	}
}

func (c CertificationResponse) Hash() string {
	return c.hash
}

func (c CertificationResponse) AnchorID() int {
	return c.anchorID
}
