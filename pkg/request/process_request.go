package request

import (
	"mime/multipart"
)

type ProcessFormIntegrityRequest struct {
	Enabled bool `form:"integrity.enabled,default=false"`
}

type ProcessFormAuthenticityRequest struct {
	Enabled   bool   `form:"authenticity.enabled,default=false"`
	KeySource string `form:"authenticity.keySource"`
	Key       string `form:"authenticity.key"`
}

type ProcessFormEncryptionRequest struct {
	Enabled   bool   `form:"encryption.enabled,default=false"`
	KeySource string `form:"encryption.keySource"`
	Key       string `form:"encryption.key"`
}

type ProcessFormAvailabilityRequest struct {
	Enabled bool   `form:"availability.enabled,default=false"`
	Type    string `form:"availability.type"`
}

type AuthenticityMetadataRequest struct {
	Enabled    bool `form:"authenticity_metadata.enabled,omitempty,default=false"`
	Signatures SignatureMetadataRequest
}

type EncryptionMetadataRequest struct {
	Enabled bool   `form:"encryption_metadata.enabled,omitempty,default=false"`
	Key     string `form:"encryption_metadata.key,omitempty"`
	Alg     string `form:"encryption_metadata.alg,omitempty"`
	Subject string `form:"encryption_metadata.subject,omitempty"`
}

type SignatureMetadataRequest struct {
	Signature   string `form:"authenticity_metadata.signature.signature,omitempty"`
	Alg         string `form:"authenticity_metadata.signature.alg,omitempty"`
	Kid         string `form:"authenticity_metadata.signature.kid,omitempty"`
	MessageHash string `form:"authenticity_metadata.signature.message_hash,omitempty"`
	Subject     string `form:"authenticity_metadata.signature.subject,omitempty"`
}

type ProcessFormRequest struct {
	File                 *multipart.FileHeader `form:"file"`
	Url                  string                `form:"url"`
	Integrity            ProcessFormIntegrityRequest
	Authenticity         ProcessFormAuthenticityRequest
	Encryption           ProcessFormEncryptionRequest
	Availability         ProcessFormAvailabilityRequest
	AuthenticityMetadata *AuthenticityMetadataRequest
	EncryptionMetadata   *EncryptionMetadataRequest
}
