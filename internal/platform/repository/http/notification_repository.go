package http

import (
	"bytes"
	"errors"
	"github.com/rs/zerolog"
	"mime/multipart"
	"net/http"
)

var ErrNotification = errors.New("notification couldn't send")

type HttpNotificationRepository struct {
	httpClient        http.Client
	clientEndpointURL string
	logger            zerolog.Logger
}

func NewHttpNotificationRepository(httpClient http.Client, clientEndpointURL string, logger zerolog.Logger) *HttpNotificationRepository {
	return &HttpNotificationRepository{httpClient: httpClient, clientEndpointURL: clientEndpointURL, logger: logger}
}

func (h HttpNotificationRepository) NotifyCertification(hash string, file []byte) error {
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
