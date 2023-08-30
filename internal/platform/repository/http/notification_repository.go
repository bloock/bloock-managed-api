package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"mime/multipart"
	"net/http"
)

var ErrNotification = errors.New("notification couldn't send")

type HttpNotificationRepository struct {
	httpClient     http.Client
	destinationURL string
	logger         zerolog.Logger
}

func NewHttpNotificationRepository(httpClient http.Client, destinationURL string, logger zerolog.Logger) *HttpNotificationRepository {
	return &HttpNotificationRepository{httpClient: httpClient, destinationURL: destinationURL, logger: logger}
}

func (h HttpNotificationRepository) NotifyCertification(hash string, whResponse interface{}, file []byte) error {
	notificationJsonBody := NotificationJsonBody{
		Hash:       hash,
		WhResponse: whResponse,
	}
	bodyBytes, err := json.Marshal(&notificationJsonBody)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		return err
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	whResponsePart, err := writer.CreateFormFile("wh_response", "wh_response.json")
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		err := ErrNotification
		return err
	}
	_, err = whResponsePart.Write(bodyBytes)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		err := ErrNotification
		return err
	}
	filePart, err := writer.CreateFormFile("file", hash)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		err := ErrNotification

		return err
	}
	_, err = filePart.Write(file)
	if err != nil {
		h.logger.Error().Err(err).Msg("")
		err := ErrNotification
		return err
	}

	resp, err := h.httpClient.Post(h.destinationURL, writer.FormDataContentType(), buf)
	if err != nil {
		err := ErrNotification
		h.logger.Error().Err(err).Msgf("response was: %s", resp.Status)
		return err
	}

	if resp.StatusCode >= 400 {
		err := ErrNotification
		h.logger.Error().Err(err).Msgf("response was: %s", resp.Status)
		return err
	}

	return nil
}

type NotificationJsonBody struct {
	Hash string `json:"hash"`
	File []byte `json:"file
"`
	WhResponse interface{} `json:"webhook_response"`
}
