package handler

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func PostProcess(processService service.BaseProcessService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var file []byte
		var isIntegrityEnabled = "false"
		var isAuthenticityEnabled = "false"
		var authenticityKeyType string
		var keyType string
		var authenticityKeyID string
		var availabilityType = "NONE"
		var useEnsResolution = "false"

		mr, err := ctx.Request.MultipartReader()
		multipartIsEmpty := err != nil
		if multipartIsEmpty {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		form, err := mr.ReadForm(config.Configuration.MaxMemory)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}
		headers := form.File["file"]
		if len(headers) == 0 {
			badRequestAPIError := NewBadRequestAPIError("file must be provided")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}
		fileReader, err := headers[0].Open()
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}
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

		isIntegrityEnabled = form.Value["integrity.enabled"][0]

		isAuthenticityEnabled = form.Value["authenticity.enabled"][0]
		if isAuthenticityEnabled == "true" {
			authenticityKeyType = form.Value["authenticity.keyType"][0]
			if authenticityKeyType == "" {
				badRequestAPIError := NewBadRequestAPIError("authenticity.keyType must be provided")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}

			keyType = form.Value["authenticity.kty"][0]
			if keyType == "" {
				badRequestAPIError := NewBadRequestAPIError("authenticity.kty must be provided")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}

			authenticityKeyID = form.Value["authenticity.key"][0]
			if authenticityKeyID == "" {
				badRequestAPIError := NewBadRequestAPIError("authenticity.key must be provided")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
			useEnsResolution = form.Value["authenticity.useEnsResolution"][0]
		}

		availabilityType = form.Value["availability.type"][0]

		signRequest, err := request.NewProcessRequest(file, isIntegrityEnabled, isAuthenticityEnabled, authenticityKeyType, keyType, authenticityKeyID, availabilityType, useEnsResolution)

		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		processResponse, err := processService.Process(ctx, *signRequest)
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusAccepted, toProcessJsonResponse(processResponse))
	}
}

func toProcessJsonResponse(processResponse *response.ProcessResponse) ProcessResponse {
	return ProcessResponse{
		Integrity: IntegrityJSONResponse{
			Hash:     processResponse.CertificationResponse().Hash(),
			AnchorId: processResponse.CertificationResponse().AnchorID(),
		},
		Authenticity: AuthenticityJSONResponse{processResponse.SignResponse().Signature()},
		Availability: AvailabilityJSONResponse{processResponse.AvailabilityResponse()},
	}
}

type ProcessResponse struct {
	Integrity    IntegrityJSONResponse    `json:"integrity"`
	Authenticity AuthenticityJSONResponse `json:"authenticity"`
	Availability AvailabilityJSONResponse `json:"availability"`
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
