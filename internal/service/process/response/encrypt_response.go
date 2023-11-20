package response

type EncryptResponse struct {
	key string
}

func NewEncryptResponse(key string) *EncryptResponse {
	return &EncryptResponse{key: key}
}

func (s EncryptResponse) Key() string {
	return s.key
}
