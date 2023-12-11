package response

import "github.com/bloock/bloock-managed-api/internal/service/process/response"

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
	Enabled    bool                 `json:"enabled"`
	Signatures []response.Signature `json:"signatures"`
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
