package handler

import (
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type postProcessForm struct {
	File                  *multipart.FileHeader `form:"file"`
	Url                   string                `form:"url"`
	IntegrityEnabled      bool                  `form:"integrity.enabled,default=false"`
	AuthenticityEnabled   bool                  `form:"authenticity.enabled,default=false"`
	AuthenticityKeySource string                `form:"authenticity.keySource"`
	AuthenticityKeyType   string                `form:"authenticity.keyType"`
	AuthenticityKey       string                `form:"authenticity.key"`
	AuthenticityUseEns    bool                  `form:"authenticity.useEnsResolution,default=false"`
	AvailabilityType      string                `form:"availability.type,default=NONE"`
	EncryptionEnabled     bool                  `form:"encryption.enabled,default=false"`
	EncryptionKeySource   string                `form:"encryption.keySource"`
	EncryptionKeyType     string                `form:"encryption.keyType"`
	EncryptionKey         string                `form:"encryption.key"`
}

func PostProcess(processService service.ProcessService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var formData postProcessForm
		err := ctx.Bind(&formData)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError("error binding form")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		var file []byte
		var filename string
		contentType := ""
		var inputUrl string
		if formData.File != nil {
			fileReader, err := formData.File.Open()
			if err != nil {
				badRequestAPIError := NewBadRequestAPIError(err.Error())
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}

			filename = formData.File.Filename
			file, err = io.ReadAll(fileReader)
			if err != nil {
				serverAPIError := NewInternalServerAPIError(err.Error())
				ctx.JSON(serverAPIError.Status, serverAPIError)
				return
			}
			if len(file) == 0 {
				badRequestAPIError := NewBadRequestAPIError("file must be a valid file")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}

			contentType = http.DetectContentType(file)
		} else if formData.Url != "" {
			u, err := url.ParseRequestURI(formData.Url)
			if err != nil {
				badRequestAPIError := NewBadRequestAPIError("Invalid URL provided")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}

			base := path.Base(u.Path)
			if base == "" {
				pathParts := strings.Split(u.Path, "/")

				// If it's empty, use the second-to-last part as the filename
				if len(pathParts) >= 2 {
					base = pathParts[len(pathParts)-2]
				}
			}
			filename = base
			inputUrl = u.String()
		} else {
			badRequestAPIError := NewBadRequestAPIError("You must provide a file or URL")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		processRequest, err := request.NewProcessRequest(file, filename, contentType, inputUrl, formData.IntegrityEnabled, formData.AuthenticityEnabled, formData.AuthenticityKeySource, formData.AuthenticityKeyType, formData.AuthenticityKey, formData.AuthenticityUseEns, formData.EncryptionEnabled, formData.EncryptionKeySource, formData.EncryptionKeyType, formData.EncryptionKey, formData.AvailabilityType)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		processResponse, err := processService.Process(ctx, *processRequest)
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusAccepted, toProcessJsonResponse(processResponse))
	}
}

func toProcessJsonResponse(processResponse *response.ProcessResponse) ProcessResponse {
	resp := ProcessResponse{
		Success: true,
		Hash:    processResponse.Hash(),
	}

	if processResponse.CertificationResponse() != nil {
		resp.Integrity = &IntegrityJSONResponse{
			AnchorId: processResponse.CertificationResponse().AnchorID(),
		}
	}

	if processResponse.SignResponse() != nil {
		resp.Authenticity = &AuthenticityJSONResponse{processResponse.SignResponse().Signature()}
	}

	if processResponse.AvailabilityResponse() != nil {
		resp.Availability = &AvailabilityJSONResponse{
			processResponse.AvailabilityResponse().Id(),
			processResponse.AvailabilityResponse().Url(),
		}
	}

	return resp
}

type ProcessResponse struct {
	Success      bool                      `json:"success"`
	Hash         string                    `json:"hash"`
	Integrity    *IntegrityJSONResponse    `json:"integrity,omitempty"`
	Authenticity *AuthenticityJSONResponse `json:"authenticity,omitempty"`
	Availability *AvailabilityJSONResponse `json:"availability,omitempty"`
}

type IntegrityJSONResponse struct {
	AnchorId int `json:"anchor_id"`
}

type AuthenticityJSONResponse struct {
	Signature string `json:"signature"`
}

type AvailabilityJSONResponse struct {
	ID  string `json:"id"`
	Url string `json:"url,omitempty"`
}
