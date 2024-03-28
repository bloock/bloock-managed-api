package response

type ProcessResponse struct {
	Success      bool                      `json:"success"`
	ProcessID    string                    `json:"process_id"`
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
	Enabled    bool                                `json:"enabled"`
	Signatures []AuthenticitySignatureJSONResponse `json:"signatures"`
}

type AuthenticitySignatureJSONResponse struct {
	Signature   string `json:"signature"`
	Alg         string `json:"alg"`
	Kid         string `json:"kid"`
	MessageHash string `json:"message_hash"`
	Subject     string `json:"subject,omitempty"`
}

type EncryptionJSONResponse struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
	Alg     string `json:"alg"`
	Subject string `json:"subject"`
}

type AvailabilityJSONResponse struct {
	Enabled     bool   `json:"enabled"`
	Type        string `json:"type"`
	ID          string `json:"id"`
	Url         string `json:"url,omitempty"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}
