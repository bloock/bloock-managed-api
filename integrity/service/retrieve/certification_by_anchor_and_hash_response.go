package retrieve

import "github.com/bloock/bloock-sdk-go/v2/entity/integrity"

type CertificationByAnchorAndHashResponse struct {
}

func NewCertificationByAnchorAndHashResponse(string, *integrity.Anchor) *CertificationByAnchorAndHashResponse {
	return &CertificationByAnchorAndHashResponse{}
}
