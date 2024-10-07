package process

import (
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"github.com/bloock/bloock-managed-api/internal/service/process"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type ListProcessResponse struct {
	utils.Pagination
	Processes []domain.Process `json:"processes"`
}

func ListProcess(l zerolog.Logger, ent *connection.EntConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pq, err := utils.NewPaginationQuery(ctx)
		if err != nil {
			pq = utils.PaginationQuery{
				Page:    1,
				PerPage: 10,
			}
		}

		service := process.NewListProcess(ctx, l, ent)
		listProcess, pagination, err := service.List(ctx, pq)
		if err != nil {
			serverAPIError := api_error.NewInternalServerAPIError(err)
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusOK, ListProcessResponse{
			Processes:  listProcess,
			Pagination: pagination,
		})
	}
}
