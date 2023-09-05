package sql

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/repository/sql/ent/certification"
	"context"
	"time"

	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/rs/zerolog"
)

type SQLCertificationRepository struct {
	connection connection.EntConnection
	dbTimeout  time.Duration
	logger     zerolog.Logger
}

func NewSQLCertificationRepository(connection connection.EntConnection, dbTimeout time.Duration, logger zerolog.Logger) *SQLCertificationRepository {
	return &SQLCertificationRepository{connection: connection, dbTimeout: dbTimeout, logger: logger}
}

func (s SQLCertificationRepository) SaveCertification(ctx context.Context, certification domain.Certification) error {
	crt := s.connection.DB().
		Certification.Create().
		SetHash(certification.Hash()).
		SetAnchorID(certification.AnchorID()).
		SetDataID(certification.DataID()).
		SetAnchor(certification.Anchor())

	if _, err := crt.Save(ctx); err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (s SQLCertificationRepository) GetCertificationsByAnchorID(ctx context.Context, anchor int) ([]domain.Certification, error) {
	certificationsSchema, err := s.connection.DB().Certification.Query().
		Where(certification.AnchorID(anchor)).All(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return []domain.Certification{}, err
	}

	var certifications []domain.Certification
	for _, crt := range certificationsSchema {
		newCrt := domain.NewCertification(crt.AnchorID, crt.Hash, crt.Anchor)
		certifications = append(certifications, *newCrt)
	}

	return certifications, nil
}

func (s SQLCertificationRepository) UpdateCertificationAnchor(ctx context.Context, anchor integrity.Anchor) error {
	if _, err := s.connection.DB().Certification.Update().SetAnchor(&anchor).
		Where(certification.AnchorIDLTE(int(anchor.Id))).
		Save(ctx); err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (s SQLCertificationRepository) UpdateCertificationDataID(ctx context.Context, hash string, dataID string) error {
	if _, err := s.connection.DB().Certification.Update().SetDataID(dataID).
		Where(certification.HashEQ(hash)).Save(ctx); err != nil {
		return err
	}

	return nil
}
