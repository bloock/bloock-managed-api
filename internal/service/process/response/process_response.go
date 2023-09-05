package response

import (
	autenticityResponse "bloock-managed-api/internal/service/authenticity/response"
	availabilityResponse "bloock-managed-api/internal/service/availability/response"
	integrityResponse "bloock-managed-api/internal/service/integrity/response"
)

type ProcessResponse struct {
	certificationResponse *integrityResponse.CertificationResponse
	signResponse          *autenticityResponse.SignResponse
	availabilityResponse  *availabilityResponse.AvailabilityResponse
}

type ProcessResponseBuilder struct {
	processResponse *ProcessResponse
}

func NewProcessResponseBuilder() *ProcessResponseBuilder {
	processResponse := &ProcessResponse{}
	b := &ProcessResponseBuilder{processResponse: processResponse}
	return b
}

func (b *ProcessResponseBuilder) CertificationResponse(certificationResponse integrityResponse.CertificationResponse) *ProcessResponseBuilder {
	b.processResponse.certificationResponse = &certificationResponse
	return b
}

func (b *ProcessResponseBuilder) SignResponse(signResponse autenticityResponse.SignResponse) *ProcessResponseBuilder {
	b.processResponse.signResponse = &signResponse
	return b
}

func (b *ProcessResponseBuilder) AvailabilityResponse(availabilityResponse availabilityResponse.AvailabilityResponse) *ProcessResponse {
	b.processResponse.availabilityResponse = &availabilityResponse
	return b.processResponse
}

func (b *ProcessResponseBuilder) Build() *ProcessResponse {
	return b.processResponse
}

func (b *ProcessResponseBuilder) CertificationHash() string {
	return b.processResponse.certificationResponse.Hash()
}

func (p ProcessResponse) CertificationResponse() *integrityResponse.CertificationResponse {
	return p.certificationResponse
}

func (p ProcessResponse) SignResponse() *autenticityResponse.SignResponse {
	return p.signResponse
}

func (p ProcessResponse) AvailabilityResponse() *availabilityResponse.AvailabilityResponse {
	return p.availabilityResponse
}
