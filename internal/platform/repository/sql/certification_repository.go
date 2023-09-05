package sql

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/repository/sql/ent/certification"
	"context"
	"time"

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
		SetHash(certification.Hash).
		SetAnchorID(certification.AnchorID).
		SetDataID(certification.DataID)

	if _, err := crt.Save(ctx); err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (s SQLCertificationRepository) GetCertificationsByAnchorID(ctx context.Context, anchorID int) ([]domain.Certification, error) {
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
		}
		certifications = append(certifications, newCrt)
	}

	return certifications, nil
}

func (s SQLCertificationRepository) ExistCertificationByHash(ctx context.Context, hash string) (bool, error) {
	exist, err := s.connection.DB().Certification.Query().
		Where(certification.HashEQ(hash)).Exist(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return false, err
	}

	return exist, nil
}

func (s SQLCertificationRepository) UpdateCertificationDataID(ctx context.Context, cert domain.Certification) error {
	if _, err := s.connection.DB().Certification.Update().
		SetDataID(cert.DataID).
		SetAnchorID(cert.AnchorID).
		Where(certification.HashEQ(cert.Hash)).Save(ctx); err != nil {
		return err
	}

	return nil
}
