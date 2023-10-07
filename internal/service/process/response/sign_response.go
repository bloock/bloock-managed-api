package response

type SignResponse struct {
	signature string
}

func NewSignResponse(signature string) *SignResponse {
	return &SignResponse{signature: signature}
}

func (s SignResponse) Signature() string {
	return s.signature
}
