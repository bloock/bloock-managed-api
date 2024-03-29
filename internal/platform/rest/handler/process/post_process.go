package process

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/bloock/bloock-managed-api/internal/domain"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/service/process"
	"github.com/bloock/bloock-managed-api/internal/service/process/request"
	"github.com/bloock/bloock-managed-api/internal/service/process/response"
	http_request "github.com/bloock/bloock-managed-api/pkg/request"
	http_response "github.com/bloock/bloock-managed-api/pkg/response"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

func PostProcess(l zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		multiPartForm, err := ctx.MultipartForm()
		if err != nil {
			badRequestAPIError := api_error.NewBadRequestAPIError("error getting form")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		var formData http_request.ProcessFormRequest
		err = ctx.Bind(&formData)
		if err != nil {
			badRequestAPIError := api_error.NewBadRequestAPIError("error binding form")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		processService := process.NewProcessService(ctx, l)

		var file domain.File
		if multiPartForm.File["file"] != nil {
			file, err = loadFile(multiPartForm.File["file"])
			if err != nil {
				badRequestAPIError := api_error.NewBadRequestAPIError(err.Error())
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
		} else if formData.Url != "" {
			u, err := url.ParseRequestURI(formData.Url)
			if err != nil {
				badRequestAPIError := api_error.NewBadRequestAPIError("Invalid URL provided")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}

			file, err = processService.LoadUrl(ctx, u)
			if err != nil {
				badRequestAPIError := api_error.NewBadRequestAPIError("Invalid URL provided")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
		} else {
			badRequestAPIError := api_error.NewBadRequestAPIError("You must provide a file or URL")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		processRequest, err := request.NewProcessRequest(file, &formData)
		if err != nil {
			badRequestAPIError := api_error.NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		processResponse, err := processService.Process(ctx, *processRequest)
		if err != nil {
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusOK, toProcessJsonResponse(processResponse))
	}
}

func toProcessJsonResponse(processResponse *response.ProcessResponse) http_response.ProcessResponse {
	resp := http_response.ProcessResponse{
		Success: true,
		Hash:    processResponse.Hash(),
	}

	if processResponse.CertificationResponse() != nil {
		resp.Integrity = &http_response.IntegrityJSONResponse{
			Enabled:  true,
			AnchorId: processResponse.CertificationResponse().AnchorID(),
		}
	}

	if processResponse.SignResponse() != nil {
		signaturesResponse := make([]http_response.AuthenticitySignatureJSONResponse, 0)
		for _, sig := range processResponse.SignResponse().Signatures() {
			signaturesResponse = append(signaturesResponse, http_response.AuthenticitySignatureJSONResponse{
				Signature:   sig.Signature,
				Alg:         sig.Alg,
				Kid:         sig.Kid,
				MessageHash: sig.MessageHash,
				Subject:     sig.Subject,
			})
		}
		resp.Authenticity = &http_response.AuthenticityJSONResponse{
			Enabled:    true,
			Signatures: signaturesResponse,
		}
	}

	if processResponse.EncryptResponse() != nil {
		resp.Encryption = &http_response.EncryptionJSONResponse{
			Enabled: true,
			Key:     processResponse.EncryptResponse().Key(),
			Alg:     processResponse.EncryptResponse().Alg(),
			Subject: processResponse.EncryptResponse().Subject(),
		}
	}

	if processResponse.AvailabilityResponse() != nil {
		resp.Availability = &http_response.AvailabilityJSONResponse{
			Enabled:     true,
			Type:        processResponse.AvailabilityResponse().Type(),
			ID:          processResponse.AvailabilityResponse().Id(),
			Url:         processResponse.AvailabilityResponse().Url(),
			ContentType: processResponse.AvailabilityResponse().ContentType(),
			Size:        processResponse.AvailabilityResponse().Size(),
		}
	}

	return resp
}

func loadFile(formsData []*multipart.FileHeader) (domain.File, error) {
	for _, formData := range formsData {
		fileReader, err := formData.Open()
		if err != nil {
			return domain.File{}, err
		}

		filename := formData.Filename
		file, err := io.ReadAll(fileReader)
		if err != nil {
			return domain.File{}, err
		}
		if len(file) == 0 {
			return domain.File{}, fmt.Errorf("file must be a valid file")
		}

		contentType := http.DetectContentType(file)

		return domain.NewFile(file, filename, contentType), nil
	}
	return domain.File{}, nil
}
