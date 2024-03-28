package response

import http_response "github.com/bloock/bloock-managed-api/pkg/response"

type ProcessResponse struct {
	hash                  string
	processID             string
	certificationResponse *IntegrityResponse
	signResponse          *SignResponse
	encryptResponse       *EncryptResponse
	availabilityResponse  *AvailabilityResponse
}

type ProcessResponseBuilder struct {
	processResponse *ProcessResponse
}

func NewProcessResponseBuilder() *ProcessResponseBuilder {
	processResponse := &ProcessResponse{}
	b := &ProcessResponseBuilder{processResponse: processResponse}
	return b
}

func (b *ProcessResponseBuilder) HashResponse(hash string) *ProcessResponseBuilder {
	b.processResponse.hash = hash
	return b
}

func (b *ProcessResponseBuilder) ProcessIDResponse(processID string) *ProcessResponseBuilder {
	b.processResponse.processID = processID
	return b
}

func (b *ProcessResponseBuilder) CertificationResponse(certificationResponse IntegrityResponse) *ProcessResponseBuilder {
	b.processResponse.certificationResponse = &certificationResponse
	return b
}

func (b *ProcessResponseBuilder) SignResponse(signResponse SignResponse) *ProcessResponseBuilder {
	b.processResponse.signResponse = &signResponse
	return b
}

func (b *ProcessResponseBuilder) EncryptResponse(encryptResponse EncryptResponse) *ProcessResponseBuilder {
	b.processResponse.encryptResponse = &encryptResponse
	return b
}

func (b *ProcessResponseBuilder) AvailabilityResponse(availabilityResponse AvailabilityResponse) *ProcessResponseBuilder {
	b.processResponse.availabilityResponse = &availabilityResponse
	return b
}

func (b *ProcessResponseBuilder) Build() *ProcessResponse {
	return b.processResponse
}

func (b *ProcessResponseBuilder) CertificationHash() string {
	return b.processResponse.certificationResponse.Hash()
}

func (p ProcessResponse) Hash() string {
	return p.hash
}

func (p ProcessResponse) ProcessID() string {
	return p.processID
}

func (p ProcessResponse) CertificationResponse() *IntegrityResponse {
	return p.certificationResponse
}

func (p ProcessResponse) SignResponse() *SignResponse {
	return p.signResponse
}

func (p ProcessResponse) EncryptResponse() *EncryptResponse {
	return p.encryptResponse
}

func (p ProcessResponse) AvailabilityResponse() *AvailabilityResponse {
	return p.availabilityResponse
}

func (p ProcessResponse) MapToHandlerProcessResponse() http_response.ProcessResponse {
	resp := http_response.ProcessResponse{
		Success:   true,
		ProcessID: p.ProcessID(),
		Hash:      p.Hash(),
	}

	if p.CertificationResponse() != nil {
		resp.Integrity = &http_response.IntegrityJSONResponse{
			Enabled:  true,
			AnchorId: p.CertificationResponse().AnchorID(),
		}
	}

	if p.SignResponse() != nil {
		signaturesResponse := make([]http_response.AuthenticitySignatureJSONResponse, 0)
		for _, sig := range p.SignResponse().Signatures() {
			signaturesResponse = append(signaturesResponse, http_response.AuthenticitySignatureJSONResponse{
				Signature:   sig.Signature,
				Alg:         sig.Alg,
				Kid:         sig.Kid,
				MessageHash: sig.MessageHash,
				Subject:     sig.Subject,
			})
		}
		resp.Authenticity = &http_response.AuthenticityJSONResponse{
			Enabled:    true,
			Signatures: signaturesResponse,
		}
	}

	if p.EncryptResponse() != nil {
		resp.Encryption = &http_response.EncryptionJSONResponse{
			Enabled: true,
			Key:     p.EncryptResponse().Key(),
			Alg:     p.EncryptResponse().Alg(),
			Subject: p.EncryptResponse().Subject(),
		}
	}

	if p.AvailabilityResponse() != nil {
		resp.Availability = &http_response.AvailabilityJSONResponse{
			Enabled:     true,
			Type:        p.AvailabilityResponse().Type(),
			ID:          p.AvailabilityResponse().Id(),
			Url:         p.AvailabilityResponse().Url(),
			ContentType: p.AvailabilityResponse().ContentType(),
			Size:        p.AvailabilityResponse().Size(),
		}
	}

	return resp
}
