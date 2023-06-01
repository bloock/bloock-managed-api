package sql

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/repository/sql/ent"
	"bloock-managed-api/internal/platform/repository/sql/ent/certification"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/rs/zerolog"
	"time"
)

type SQLCertificationRepository struct {
	connection connection.EntConnection
	dbTimeout  time.Duration
	logger     zerolog.Logger
}

func NewSQLCertificationRepository(connection connection.EntConnection, dbTimeout time.Duration, logger zerolog.Logger) *SQLCertificationRepository {
	return &SQLCertificationRepository{connection: connection, dbTimeout: dbTimeout, logger: logger}
}

func (s SQLCertificationRepository) SaveCertification(ctx context.Context, certifications []domain.Certification) error {
	var certificationsCreate []*ent.CertificationCreate
	for _, crt := range certifications {
		certificationsCreate = append(certificationsCreate, s.connection.DB().
			Certification.Create().
			SetHash(crt.Hash()).
			SetAnchorID(crt.AnchorID()).
			SetAnchor(crt.Anchor()))
	}

	if _, err := s.connection.DB().Certification.CreateBulk(certificationsCreate...).Save(ctx); err != nil {
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
