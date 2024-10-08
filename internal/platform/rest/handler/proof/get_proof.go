package proof

import (
	"errors"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/service/proof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type GetProofRequest struct {
	Message []string `json:"messages" binding:"required"`
}

type GetProofResponse struct {
	Leaves []string    `json:"leaves"`
	Nodes  []string    `json:"nodes"`
	Bitmap string      `json:"bitmap"`
	Depth  string      `json:"depth"`
	Root   string      `json:"root"`
	Anchor interface{} `json:"anchor"`
}

func mapToGetProofResponse(proof domain.BloockProof) GetProofResponse {
	return GetProofResponse{
		Leaves: proof.Leaves,
		Nodes:  proof.Nodes,
		Bitmap: proof.Bitmap,
		Depth:  proof.Depth,
		Root:   proof.Root,
		Anchor: proof.Anchor,
	}
}

func GetProof(l zerolog.Logger, ent *connection.EntConnection, maxProofMessageSize int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetProofRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errBind := api_error.NewBadRequestAPIError(err.Error())
			ctx.JSON(errBind.Status, err.Error())
			return
		}

		service := proof.NewGetProof(ctx, l, ent, maxProofMessageSize)

		res, err := service.Get(ctx, req.Message)
		if err != nil {
			if errors.Is(proof.ErrMessageNotFound, err) || errors.Is(proof.ErrEmptyMessages, err) || errors.Is(proof.ErrInvalidMessageHash, err) ||
				errors.Is(proof.ErrMaxProofMessagesSize, err) || errors.Is(proof.ErrInconsistentMessages, err) {
				notFoundAPIError := api_error.NewAPIError(http.StatusNotFound, err.Error())
				ctx.JSON(notFoundAPIError.Status, notFoundAPIError)
				return
			}
			if errors.Is(repository.ErrUnreadyProofStatus, err) {
				badRequestAPIError := api_error.NewAPIError(http.StatusBadRequest, err.Error())
				ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
				return
			}
			serverAPIError := api_error.NewInternalServerAPIError(err)
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusOK, mapToGetProofResponse(res))
	}
}
