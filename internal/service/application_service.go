package service

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/domain"
	processRequest "github.com/bloock/bloock-managed-api/internal/service/process/request"
	"github.com/bloock/bloock-managed-api/internal/service/process/response"
)

type ProcessService interface {
	LoadUrl(ctx context.Context, url string) ([]domain.File, error)
	Process(ctx context.Context, req processRequest.ProcessRequest) (*response.ProcessResponse, error)
}

type NotifyService interface {
	Notify(ctx context.Context, anchorID int) error
}
