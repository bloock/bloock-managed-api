package response

type ProcessResponse struct {
	Success      bool                      `json:"success"`
	Hash         string                    `json:"hash"`
	Integrity    *IntegrityJSONResponse    `json:"integrity,omitempty"`
	Authenticity *AuthenticityJSONResponse `json:"authenticity,omitempty"`
	Encryption   *EncryptionJSONResponse   `json:"encryption,omitempty"`
	Availability *AvailabilityJSONResponse `json:"availability,omitempty"`
}

type IntegrityJSONResponse struct {
	AnchorId int `json:"anchor_id"`
}

type AuthenticityJSONResponse struct {
	Key       string `json:"key"`
	Signature string `json:"signature"`
}

type EncryptionJSONResponse struct {
	Key       string `json:"key"`
	Signature string `json:"signature"`
}

type AvailabilityJSONResponse struct {
	ID  string `json:"id"`
	Url string `json:"url,omitempty"`
}
