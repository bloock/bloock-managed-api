package request

type ProcessFormIntegrityRequest struct {
	Enabled   bool `form:"integrity.enabled,default=false"`
	Aggregate bool `form:"integrity.aggregate,default=false"`
}

type ProcessFormAuthenticityRequest struct {
	Enabled       bool   `form:"authenticity.enabled,default=false"`
	KeySource     string `form:"authenticity.keySource,omitempty"`
	Key           string `form:"authenticity.key,omitempty"`
	AccessEnabled bool   `form:"authenticity.accessEnabled,default=false"`
	AccessType    string `form:"authenticity.accessType,omitempty"`
	AccessCode    string `form:"authenticity.accessCode,omitempty"`
}

type ProcessFormEncryptionRequest struct {
	Enabled       bool   `form:"encryption.enabled,default=false"`
	KeySource     string `form:"encryption.keySource,omitempty"`
	Key           string `form:"encryption.key,omitempty"`
	AccessEnabled bool   `form:"encryption.accessEnabled,default=false"`
	AccessType    string `form:"encryption.accessType,omitempty"`
	AccessCode    string `form:"encryption.accessCode,omitempty"`
}

type ProcessFormAvailabilityRequest struct {
	Enabled bool   `form:"availability.enabled,default=false"`
	Type    string `form:"availability.type,omitempty"`
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
	Url          string `form:"url,omitempty"`
	Integrity    ProcessFormIntegrityRequest
	Authenticity ProcessFormAuthenticityRequest
	Encryption   ProcessFormEncryptionRequest
	Availability ProcessFormAvailabilityRequest
}
