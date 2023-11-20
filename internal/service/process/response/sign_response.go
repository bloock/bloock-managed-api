package response

type SignResponse struct {
	key       string
	signature string
}

func NewSignResponse(key, signature string) *SignResponse {
	return &SignResponse{key: key, signature: signature}
}

func (s SignResponse) Signature() string {
	return s.signature
}

func (s SignResponse) Key() string {
	return s.key
}
