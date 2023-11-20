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
	Enabled  bool `json:"enabled"`
	AnchorId int  `json:"anchor_id"`
}

type AuthenticityJSONResponse struct {
	Enabled   bool   `json:"enabled"`
	Key       string `json:"key"`
	Signature string `json:"signature"`
}

type EncryptionJSONResponse struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
}

type AvailabilityJSONResponse struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	ID      string `json:"id"`
	Url     string `json:"url,omitempty"`
}
