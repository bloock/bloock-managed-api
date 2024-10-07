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

func GetFileByHash(l zerolog.Logger, ent *connection.EntConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hash := ctx.Param("id")
		if hash == "" {
			badRequestAPIError := api_error.NewAPIError(http.StatusNotFound, "empty hash id")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		getHashService := process.NewGetFileByHash(ctx, l, ent)

		fileBytes, err := getHashService.Get(ctx, hash)
		if err != nil {
			if errors.Is(process.ErrHashNotFound, err) {
				notFoundAPIError := api_error.NewAPIError(http.StatusNotFound, err.Error())
				ctx.JSON(notFoundAPIError.Status, notFoundAPIError)
				return
			}
			serverAPIError := api_error.NewInternalServerAPIError(err)
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		if _, err = ctx.Writer.Write(fileBytes); err != nil {
			serverAPIError := api_error.NewInternalServerAPIError(err)
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}
	}
}
