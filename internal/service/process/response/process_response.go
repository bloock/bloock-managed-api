package response

type ProcessResponse struct {
	hash                  string
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
