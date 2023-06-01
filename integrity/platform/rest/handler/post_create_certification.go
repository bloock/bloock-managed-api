package handler

import (
	"bloock-managed-api/integrity/service/create"
	"bloock-managed-api/integrity/service/response"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func PostCreateCertification(certification create.Certification) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		mr, err := ctx.Request.MultipartReader()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		hasJsonBody := mr == nil
		var files [][]byte
		var file []byte

		if hasJsonBody {
			if _, err := ctx.Request.Body.Read(file); err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			for true {
				p, err := mr.NextPart()
				if errors.Is(err, io.EOF) {
					break
				}

				if _, err := p.Read(file); err != nil {
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
