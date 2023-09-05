package handler

import (
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type postProcessForm struct {
	File                  *multipart.FileHeader `form:"file" binding:"required"`
	IntegrityEnabled      bool                  `form:"integrity.enabled,default=false"`
	AuthenticityEnabled   bool                  `form:"authenticity.enabled,default=false"`
	AuthenticityKeySource string                `form:"authenticity.keySource"`
	AuthenticityKeyType   string                `form:"authenticity.keyType"`
	AuthenticityKey       string                `form:"authenticity.key"`
	AuthenticityUseEns    bool                  `form:"authenticity.useEnsResolution,default=false"`
	AvailabilityType      string                `form:"availability.type,default=NONE"`
}

func PostProcess(processService service.BaseProcessService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var formData postProcessForm
		err := ctx.Bind(&formData)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError("error binding form")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		fileReader, err := formData.File.Open()
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}
		file, err := io.ReadAll(fileReader)
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

		processRequest, err := request.NewProcessRequest(file, formData.IntegrityEnabled, formData.AuthenticityEnabled, formData.AuthenticityKeySource, formData.AuthenticityKeyType, formData.AuthenticityKey, formData.AuthenticityUseEns, formData.AvailabilityType)
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
	response := ProcessResponse{
		Success: true,
	}

	if processResponse.CertificationResponse() != nil {
		response.Integrity = &IntegrityJSONResponse{
			Hash:     processResponse.CertificationResponse().Hash(),
			AnchorId: processResponse.CertificationResponse().AnchorID(),
		}
	}

	if processResponse.SignResponse() != nil {
		response.Authenticity = &AuthenticityJSONResponse{processResponse.SignResponse().Signature()}
	}

	if processResponse.AvailabilityResponse() != nil {
		response.Availability = &AvailabilityJSONResponse{processResponse.AvailabilityResponse().Cid()}
	}

	return response
}

type ProcessResponse struct {
	Success      bool                      `json:"success"`
	Integrity    *IntegrityJSONResponse    `json:"integrity,omitempty"`
	Authenticity *AuthenticityJSONResponse `json:"authenticity,omitempty"`
	Availability *AvailabilityJSONResponse `json:"availability,omitempty"`
}

type IntegrityJSONResponse struct {
	Hash     string `json:"hash"`
	AnchorId int    `json:"anchor_id"`
}

type AuthenticityJSONResponse struct {
	Signature string `json:"signature"`
}
type AvailabilityJSONResponse struct {
	ID string `json:"id"`
}
type CertificationJSONRequest struct {
	Data interface{}
}
