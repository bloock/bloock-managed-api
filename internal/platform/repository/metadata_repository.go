package repository

import (
	"context"
	"time"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/certification"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockMetadataRepository struct {
	recordClient client.RecordClient
	connection   *connection.EntConnection
	dbTimeout    time.Duration
	logger       zerolog.Logger
}

func NewBloockMetadataRepository(ctx context.Context, logger zerolog.Logger) repository.MetadataRepository {
	logger.With().Caller().Str("component", "metadata-repository").Logger()

	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection(config.Configuration.Db.ConnectionString, entConnector, logger)
	if err != nil {
		panic(err)
	}
	err = conn.Migrate()
	if err != nil {
		panic(err)
	}

	return &BloockMetadataRepository{
		recordClient: client.NewRecordClient(),
		connection:   conn,
		dbTimeout:    5 * time.Second,
		logger:       logger,
	}
}

func (f BloockMetadataRepository) GetRecord(ctx context.Context, file []byte) (*record.Record, error) {
	record, err := f.recordClient.FromFile(file).Build()
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (f BloockMetadataRepository) GetRecordDetails(ctx context.Context, file []byte) (*record.RecordDetails, error) {
	record, err := f.recordClient.FromFile(file).GetDetails()
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (f BloockMetadataRepository) GetFileHash(ctx context.Context, file []byte) (string, error) {
	record, err := f.recordClient.FromFile(file).Build()
	if err != nil {
		return "", err
	}
	hash, err := record.GetHash()
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (s BloockMetadataRepository) SaveCertification(ctx context.Context, certification domain.Certification) error {
	crt := s.connection.DB().
		Certification.Create().
		SetHash(certification.Hash).
		SetAnchorID(certification.AnchorID).
		SetDataID(certification.DataID)

	if _, err := crt.Save(ctx); err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (c BloockMetadataRepository) UpdateCertification(ctx context.Context, certification domain.Certification) error {
	exist, err := c.ExistCertificationByHash(ctx, certification.Hash)
	if err != nil {
		return err
	}
	if !exist {
		return c.SaveCertification(ctx, certification)
	} else {
		return c.UpdateCertificationDataID(ctx, certification)
	}
}

func (s BloockMetadataRepository) GetCertificationsByAnchorID(ctx context.Context, anchorID int) ([]domain.Certification, error) {
	certificationsSchema, err := s.connection.DB().Certification.Query().
		Where(certification.AnchorID(anchorID)).All(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return []domain.Certification{}, err
	}

	var certifications []domain.Certification
	for _, crt := range certificationsSchema {
		newCrt := domain.Certification{
			AnchorID: crt.AnchorID,
			Hash:     crt.Hash,
			DataID:   crt.DataID,
		}
		certifications = append(certifications, newCrt)
	}

	return certifications, nil
}

func (s BloockMetadataRepository) ExistCertificationByHash(ctx context.Context, hash string) (bool, error) {
	exist, err := s.connection.DB().Certification.Query().
		Where(certification.HashEQ(hash)).Exist(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return false, err
	}

	return exist, nil
}

func (s BloockMetadataRepository) UpdateCertificationDataID(ctx context.Context, cert domain.Certification) error {
	if _, err := s.connection.DB().Certification.Update().
		SetDataID(cert.DataID).
		SetAnchorID(cert.AnchorID).
		Where(certification.HashEQ(cert.Hash)).Save(ctx); err != nil {
		return err
	}

	return nil
}
