package handler

import (
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/response"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func PostCreateCertification(certification create.Certification) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		mr, err := ctx.Request.MultipartReader()

		var files [][]byte

		if err != nil {
			jsonBytes, err := io.ReadAll(ctx.Request.Body)
			files = append(files, jsonBytes)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

		} else {
			for true {
				p, err := mr.NextPart()
				if errors.Is(err, io.EOF) {
					break
				}

				file, err := io.ReadAll(p)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, err.Error())
					return
				}

				files = append(files, file)
			}
		}

		certificationResponse, err := certification.Certify(ctx, files)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		responses := mapToCertificationJsonResponse(certificationResponse)
		ctx.JSON(http.StatusAccepted, responses)
	}
}

func mapToCertificationJsonResponse(certificationResponse []response.CertificationResponse) []CertificationJSONResponse {
	var responses []CertificationJSONResponse
	for _, crt := range certificationResponse {
		response := CertificationJSONResponse{
			Hash:     crt.Hash(),
			AnchorId: crt.AnchorID(),
		}

		responses = append(responses, response)
	}
	return responses
}

type CertificationJSONResponse struct {
	Hash     string `json:"hash"`
	AnchorId int    `json:"anchor_id"`
}
