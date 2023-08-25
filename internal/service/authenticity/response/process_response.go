package response

import "bloock-managed-api/internal/service/integrity/response"

type ProcessResponse struct {
	certificationResponse []response.CertificationResponse
	signResponse          SignResponse
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

func (b *ProcessResponseBuilder) CertificationResponse(certificationResponse []response.CertificationResponse) *ProcessResponseBuilder {
	b.processResponse.certificationResponse = certificationResponse
	return b
}

func (b *ProcessResponseBuilder) SignResponse(signResponse SignResponse) *ProcessResponseBuilder {
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
