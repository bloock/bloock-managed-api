package aggregate

import (
	"errors"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/service/aggregate"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type PutAggregateResponse struct {
	Success bool `json:"success"`
}

func PutAggregate(l zerolog.Logger, ent *connection.EntConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		service := aggregate.NewServiceAggregator(ctx, l, ent)

		if err := service.Aggregate(ctx); err != nil {
			if errors.Is(err, aggregate.ErrApiKeyNotFound) || errors.Is(err, aggregate.ErrMinimumPendingMessages) {
				badAPIError := api_error.NewBadRequestAPIError(err.Error())
				ctx.JSON(badAPIError.Status, badAPIError)
				return
			}
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusOK, PutAggregateResponse{Success: true})
	}
}
