package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bloock/bloock-managed-api/internal/config"
	"io"
	"net/http"
	"strings"

	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/pkg"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

var ErrInvalidResponse = errors.New("record couldn't send")
var ErrUnreadyProofStatus = errors.New("unready proof status")

type BloockIntegrityRepository struct {
	httpClient http.Client
	client     client.BloockClient
	apiKey     string
	apiVersion string
	logger     zerolog.Logger
}

func NewBloockIntegrityRepository(ctx context.Context, l zerolog.Logger) repository.IntegrityRepository {
	logger := l.With().Caller().Str("component", "integrity-repository").Logger()

	apiKey := pkg.GetApiKeyFromContext(ctx)
	c := client.NewBloockClient(apiKey, nil)

	return &BloockIntegrityRepository{
		httpClient: http.Client{},
		apiKey:     apiKey,
		client:     c,
		apiVersion: config.Configuration.Api.ApiVersion,
		logger:     logger,
	}
}

type sendRecordRequest struct {
	Messages []string `json:"messages"`
}

type sendRecordResponse struct {
	Anchor int `json:"anchor"`
}

type getProofRequest struct {
	Messages []string `json:"messages"`
}

func mapToSendRecordRequest(hash string) sendRecordRequest {
	sr := sendRecordRequest{
		Messages: []string{hash},
	}
	return sr
}

func mapToGetProofRequest(hash []string) getProofRequest {
	pr := getProofRequest{
		Messages: hash,
	}
	return pr
}

func (b BloockIntegrityRepository) Certify(ctx context.Context, file []byte) (domain.Certification, error) {
	rec, err := client.NewRecordClient().FromFile(file).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("error certifying data")
		return domain.Certification{}, err
	}

	receipt, err := b.client.IntegrityClient.SendRecords([]record.Record{rec})
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return domain.Certification{}, err
	}
	dataHash, err := rec.GetHash()
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return domain.Certification{}, err
	}

	return domain.Certification{
		AnchorID: int(receipt[0].Anchor),
		Data:     file,
		Hash:     dataHash,
		Record:   &rec,
	}, nil
}

func (b BloockIntegrityRepository) CertifyFromHash(ctx context.Context, hash string, apiKey string) (domain.Certification, error) {
	url := fmt.Sprintf("%s/records/v1/records", config.Configuration.Bloock.ApiHost)

	auth := b.apiKey
	if apiKey != "" {
		auth = apiKey
	}
	var recordResp sendRecordResponse
	if err := b.postRequest(url, mapToSendRecordRequest(hash), &recordResp, auth); err != nil {
		return domain.Certification{}, err
	}

	return domain.Certification{
		Hash:     hash,
		AnchorID: recordResp.Anchor,
	}, nil
}

func (b BloockIntegrityRepository) GetProof(ctx context.Context, hash []string, apiKey string) (domain.BloockProof, error) {
	url := fmt.Sprintf("%s/proof/v1/proof", config.Configuration.Bloock.ApiHost)

	auth := b.apiKey
	if apiKey != "" {
		auth = apiKey
	}
	var proofResp json.RawMessage
	if err := b.postRequest(url, mapToGetProofRequest(hash), &proofResp, auth); err != nil {
		return domain.BloockProof{}, err
	}

	var bloockProof domain.BloockProof
	if err := json.Unmarshal(proofResp, &bloockProof); err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return domain.BloockProof{}, err
	}

	return bloockProof, nil
}

func (b BloockIntegrityRepository) postRequest(url string, body interface{}, response interface{}, apiKey string) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return err
	}
	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_version", b.apiVersion)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		b.logger.Error().Err(err).Msgf("response was: %s", resp.Status)
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == 404 {
			respByte, err := io.ReadAll(resp.Body)
			if err != nil {
				b.logger.Error().Err(err).Msg(err.Error())
				return err
			}
			if strings.Contains(string(respByte), "proof not found") || strings.Contains(string(respByte), "records not found") {
				err = ErrUnreadyProofStatus
				b.logger.Error().Err(err).Msg("")
				return err
			}
		}
		err = ErrInvalidResponse
		b.logger.Error().Err(err).Msgf("response was: %s", resp.Status)
		return err
	}

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return err
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
