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

func GetProcessByID(l zerolog.Logger, ent *connection.EntConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			badRequestAPIError := api_error.NewAPIError(http.StatusNotFound, "empty process id")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		processService := process.NewGetProcessByID(ctx, l, ent)

		res, err := processService.Get(ctx, id)
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
			serverAPIError := api_error.NewInternalServerAPIError(err)
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}
