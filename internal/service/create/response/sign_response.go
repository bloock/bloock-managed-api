package response

type SignResponse struct {
	record []byte
}

func (s SignResponse) Record() []byte {
	return s.record
}

func NewSignResponse(record []byte) *SignResponse {
	return &SignResponse{record: record}
}
