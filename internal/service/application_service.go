package service

import (
	processRequest "bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"context"
)

type ProcessService interface {
	Process(ctx context.Context, req processRequest.ProcessRequest) (*response.ProcessResponse, error)
}

type NotifyService interface {
	Notify(ctx context.Context, anchorID int) error
}
