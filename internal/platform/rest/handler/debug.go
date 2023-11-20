package handler

import (
	"io"
	"mime/multipart"
	"net/http"
	"regexp"

	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/gin-gonic/gin"
)

type debugRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type debugResponse struct {
	Success bool `json:"success"`
}

func Debug() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var formData debugRequest
		err := ctx.Bind(&formData)
		if err != nil {
			badRequestAPIError := api_error.NewBadRequestAPIError("error binding form")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		fileReader, err := formData.File.Open()
		if err != nil {
			badRequestAPIError := api_error.NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		file, err := io.ReadAll(fileReader)
		if err != nil {
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}
		if len(file) == 0 {
			badRequestAPIError := api_error.NewBadRequestAPIError("empty file")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		fileName := formData.File.Filename
		pattern := "^[a-fA-F0-9]{64}$"
		regex := regexp.MustCompile(pattern)
		if !regex.MatchString(fileName) {
			badRequestAPIError := api_error.NewBadRequestAPIError("invalid sha256 hash file name")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		ctx.JSON(http.StatusOK, debugResponse{Success: true})
	}
}
