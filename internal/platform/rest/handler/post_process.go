package handler

import (
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
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
		var availabilityType string
		var useEnsResolution string = "false"

		mr, err := ctx.Request.MultipartReader()
		multipartIsEmpty := err != nil
		if multipartIsEmpty {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		for {

			p, err := mr.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}

			if p.FormName() == "file" {
				file, err = io.ReadAll(p)
				if err != nil {
					serverAPIError := NewInternalServerAPIError(err.Error())
					ctx.JSON(serverAPIError.Status, serverAPIError)
					return
				}
			}

			if p.FormName() == "integrity.enabled" {
				isIntegrityEnabled = readProp(ctx, p)
			}
			if p.FormName() == "authenticity.enabled" {
				isAuthenticityEnabled = readProp(ctx, p)
			}
			if p.FormName() == "authenticity.keyType" {
				authenticityKeyType = readProp(ctx, p)
			}
			if p.FormName() == "authenticity.kty" {
				keyType = readProp(ctx, p)
			}
			if p.FormName() == "authenticity.key" {
				authenticityKeyID = readProp(ctx, p)
			}
			if p.FormName() == "availability.type" {
				availabilityType = readProp(ctx, p)
			}
			if p.FormName() == "authenticity.useEnsResolution" {
				useEnsResolution = readProp(ctx, p)
			}
		}

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

func readProp(ctx *gin.Context, p *multipart.Part) string {
	bytes, err := io.ReadAll(p)
	if err != nil {
		serverAPIError := NewInternalServerAPIError(err.Error())
		ctx.JSON(serverAPIError.Status, serverAPIError)
		return ""
	}
	return string(bytes)
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
	ID string `json:"ID"`
}
type CertificationJSONRequest struct {
	Data interface{}
}
