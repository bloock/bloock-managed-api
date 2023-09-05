package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"

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
			badRequestAPIError := NewBadRequestAPIError("empty file")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		fileName := formData.File.Filename
		fmt.Println(fileName)
		pattern := "^[a-fA-F0-9]{64}$"
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(fileName) {
			badRequestAPIError := NewBadRequestAPIError("invalid sha256 hash file name")
			fmt.Printf("3-> %+v", err)
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		ctx.JSON(http.StatusOK, debugResponse{Success: true})
	}
}
