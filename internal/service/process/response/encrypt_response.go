package response

type EncryptResponse struct {
	key     string
	alg     string
	subject string
}

func NewEncryptResponse(key *string, alg *string, subject *string) *EncryptResponse {
	var kk, al, sub string
	if key != nil {
		kk = *key
	}
	if alg != nil {
		al = *alg
	}
	if subject != nil {
		sub = *subject
	}
	return &EncryptResponse{key: kk, alg: al, subject: sub}
}

func (s EncryptResponse) Key() string {
	return s.key
}

func (s EncryptResponse) Alg() string {
	return s.alg
}

func (s EncryptResponse) Subject() string {
	return s.subject
}
