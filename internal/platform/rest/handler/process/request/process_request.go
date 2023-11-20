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

type ProcessFormRequest struct {
	File         *multipart.FileHeader `form:"file"`
	Url          string                `form:"url"`
	Integrity    ProcessFormIntegrityRequest
	Authenticity ProcessFormAuthenticityRequest
	Encryption   ProcessFormEncryptionRequest
	Availability ProcessFormAvailabilityRequest
}
