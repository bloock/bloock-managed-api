package response

type SignResponse struct {
	signatures []Signature
}

type Signature struct {
	Signature   string `json:"signature"`
	Alg         string `json:"alg"`
	Kid         string `json:"kid"`
	MessageHash string `json:"message_hash"`
	Subject     string `json:"subject,omitempty"`
}

func NewSignature(signature string, alg string, kid string, messageHash string, subj *string) Signature {
	var subject string
	if subj != nil {
		subject = *subj
	}
	return Signature{
		Signature:   signature,
		Alg:         alg,
		Kid:         kid,
		MessageHash: messageHash,
		Subject:     subject,
	}
}

func NewSignResponse(signatures []Signature) *SignResponse {
	return &SignResponse{signatures: signatures}
}

func (s SignResponse) Signatures() []Signature {
	return s.signatures
}
