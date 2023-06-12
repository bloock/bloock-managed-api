package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"net/http"
)

type HttpNotificationRepository struct {
	httpClient     http.Client
	destinationURL string
	logger         zerolog.Logger
}

func NewHttpNotificationRepository(httpClient http.Client, destinationURL string, logger zerolog.Logger) *HttpNotificationRepository {
	return &HttpNotificationRepository{httpClient: httpClient, destinationURL: destinationURL, logger: logger}
}

func (h HttpNotificationRepository) NotifyCertification(hash string, whResponse interface{}) error {
	notificationJsonBody := NotificationJsonBody{
		Hash:       hash,
		WhResponse: whResponse,
	}
	bodyBytes, err := json.Marshal(&notificationJsonBody)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		return err
	}

	response, err := h.httpClient.Post(h.destinationURL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		return err
	}

	if response.StatusCode != http.StatusOK {
		err := errors.New("notification couldn't send")
		h.logger.Error().Err(err).Msgf("response was: %s", response.Status)
		return err
	}

	return nil
}

type NotificationJsonBody struct {
	Hash       string      `json:"hash"`
	WhResponse interface{} `json:"webhook_response"`
}
