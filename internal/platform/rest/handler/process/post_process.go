package process

import (
	"errors"
	"fmt"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/bloock/bloock-managed-api/internal/domain"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/service/process"
	"github.com/bloock/bloock-managed-api/internal/service/process/request"
	http_request "github.com/bloock/bloock-managed-api/pkg/request"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

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
			if errors.Is(process.ErrAggregateModeDisabled, err) {
				badRequestAPIError := api_error.NewBadRequestAPIError(err.Error())
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusOK, processResponse.MapToHandlerProcessResponse())
	}
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
