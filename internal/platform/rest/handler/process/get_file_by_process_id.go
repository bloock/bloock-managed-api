package process

import (
	"errors"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/service/process"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

func GetFileByProcessID(l zerolog.Logger, ent *connection.EntConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		processID := ctx.Param("id")
		if processID == "" {
			badRequestAPIError := api_error.NewAPIError(http.StatusNotFound, "empty process id")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		getProcessService := process.NewGetProcessByID(ctx, l, ent)
		processResponse, err := getProcessService.Get(ctx, processID)
		if err != nil {
			if errors.Is(process.ErrInvalidUUID, err) {
				badRequestAPIError := api_error.NewBadRequestAPIError(err.Error())
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
			if errors.Is(process.ErrProcessNotFound, err) {
				notFoundAPIError := api_error.NewAPIError(http.StatusNotFound, err.Error())
				ctx.JSON(notFoundAPIError.Status, notFoundAPIError)
				return
			}
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		getHashService := process.NewGetFileByHash(ctx, l, ent)
		fileBytes, err := getHashService.Get(ctx, processResponse.Hash)
		if err != nil {
			if errors.Is(process.ErrHashNotFound, err) {
				notFoundAPIError := api_error.NewAPIError(http.StatusNotFound, err.Error())
				ctx.JSON(notFoundAPIError.Status, notFoundAPIError)
				return
			}
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		if _, err = ctx.Writer.Write(fileBytes); err != nil {
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}
	}
}
