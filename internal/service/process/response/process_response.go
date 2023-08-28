package response

import (
	response2 "bloock-managed-api/internal/service/authenticity/response"
	"bloock-managed-api/internal/service/integrity/response"
)

type ProcessResponse struct {
	certificationResponse response.CertificationResponse
	signResponse          response2.SignResponse
	availabilityResponse  string
}

type ProcessResponseBuilder struct {
	processResponse *ProcessResponse
}

func NewProcessResponseBuilder() *ProcessResponseBuilder {
	processResponse := &ProcessResponse{}
	b := &ProcessResponseBuilder{processResponse: processResponse}
	return b
}

func (b *ProcessResponseBuilder) CertificationResponse(certificationResponse response.CertificationResponse) *ProcessResponseBuilder {
	b.processResponse.certificationResponse = certificationResponse
	return b
}

func (b *ProcessResponseBuilder) SignResponse(signResponse response2.SignResponse) *ProcessResponseBuilder {
	b.processResponse.signResponse = signResponse
	return b
}

func (b *ProcessResponseBuilder) Build() *ProcessResponse {
	return b.processResponse
}

func (b *ProcessResponseBuilder) AvailabilityResponse(url string) *ProcessResponse {
	b.processResponse.availabilityResponse = url
	return b.processResponse
}

func (p ProcessResponse) CertificationResponse() response.CertificationResponse {
	return p.certificationResponse
}

func (p ProcessResponse) SignResponse() response2.SignResponse {
	return p.signResponse
}

func (p ProcessResponse) AvailabilityResponse() string {
	return p.availabilityResponse
}
