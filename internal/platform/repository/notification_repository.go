package repository

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/rs/zerolog"
)

var ErrNotification = errors.New("notification couldn't send")

type HttpNotificationRepository struct {
	httpClient        http.Client
	clientEndpointURL string
	logger            zerolog.Logger
}

func NewHttpNotificationRepository(ctx context.Context, logger zerolog.Logger) repository.NotificationRepository {
	logger.With().Caller().Str("component", "notification-repository").Logger()

	return &HttpNotificationRepository{httpClient: http.Client{}, clientEndpointURL: config.Configuration.Webhook.ClientEndpointUrl, logger: logger}
}

func (h HttpNotificationRepository) NotifyCertification(hash string, file []byte) error {
	if h.clientEndpointURL == "" {
		return nil
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	filePart, err := writer.CreateFormFile("file", hash)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		return ErrNotification
	}
	_, err = filePart.Write(file)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		return ErrNotification
	}
	header := writer.FormDataContentType()
	_ = writer.Close()

	resp, err := h.httpClient.Post(h.clientEndpointURL, header, buf)
	if err != nil {
		err = ErrNotification
		h.logger.Error().Err(err).Msgf("response was: %s", resp.Status)
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = ErrNotification
		h.logger.Error().Err(err).Msgf("response was: %s", resp.Status)
		return err
	}

	return nil
}
