package process

import (
	"errors"
	"fmt"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/service/process"
	"github.com/bloock/bloock-managed-api/internal/service/process/request"
	"github.com/bloock/bloock-managed-api/internal/service/process/response"
	http_request "github.com/bloock/bloock-managed-api/pkg/request"
	"github.com/rs/zerolog"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func mapToPostProcessResponse(resp []response.ProcessResponse) interface{} {
	if len(resp) == 1 {
		return resp[0].MapToHandlerProcessResponse()
	}
	return response.MapToHandlerArrayProcessResponse(resp)
}

func PostProcess(l zerolog.Logger, ent *connection.EntConnection) gin.HandlerFunc {
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

		processService := process.NewProcessService(ctx, l, ent)

		var files []domain.File
		if existFilesMultipart(multiPartForm.File) {
			files, err = loadFiles(multiPartForm.File)
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

			file, err := processService.LoadUrl(ctx, u)
			if err != nil {
				badRequestAPIError := api_error.NewBadRequestAPIError("Invalid URL provided")
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
			files = append(files, file)
		} else {
			badRequestAPIError := api_error.NewBadRequestAPIError("You must provide a file or URL")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		responses := make([]response.ProcessResponse, 0)
		for _, file := range files {
			processRequest, err := request.NewProcessRequest(file, formData)
			if err != nil {
				badRequestAPIError := api_error.NewBadRequestAPIError(err.Error())
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
			processResponse, err := processService.Process(ctx, *processRequest)
			if err != nil {
				if errors.Is(process.ErrAggregateModeDisabled, err) {
					badRequestAPIError := api_error.NewBadRequestAPIError(err.Error())
					ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
					return
				}
				serverAPIError := api_error.NewInternalServerAPIError(err.Error())
				ctx.JSON(serverAPIError.Status, serverAPIError)
				return
			}
			responses = append(responses, *processResponse)
		}

		ctx.JSON(http.StatusOK, mapToPostProcessResponse(responses))
	}
}

func loadFiles(files map[string][]*multipart.FileHeader) ([]domain.File, error) {
	filesDomain := make([]domain.File, 0)
	for _, formsData := range files {
		for _, formData := range formsData {
			fileReader, err := formData.Open()
			if err != nil {
				return []domain.File{}, err
			}

			filename := formData.Filename
			file, err := io.ReadAll(fileReader)
			if err != nil {
				return []domain.File{}, err
			}
			if len(file) == 0 {
				return []domain.File{}, fmt.Errorf("file must be a valid file")
			}

			contentType := http.DetectContentType(file)

			filesDomain = append(filesDomain, domain.NewFile(file, filename, contentType))
		}
	}

	return filesDomain, nil
}

func existFilesMultipart(files map[string][]*multipart.FileHeader) bool {
	ok := false
	for _, values := range files {
		if values != nil {
			ok = true
		}
	}
	return ok
}
