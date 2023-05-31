package sql

import (
	"bloock-managed-api/ent"
	"bloock-managed-api/ent/certification"
	"bloock-managed-api/integrity/domain"
	"bloock-managed-api/integrity/platform/repository/sql/connection"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/rs/zerolog"
	"time"
)

type SQLCertificationRepository struct {
	connection connection.Connection
	dbTimeout  time.Duration
	logger     zerolog.Logger
}

func NewSQLCertificationRepository(connection connection.Connection, dbTimeout time.Duration, logger zerolog.Logger) *SQLCertificationRepository {
	return &SQLCertificationRepository{connection: connection, dbTimeout: dbTimeout, logger: logger}
}

func (s SQLCertificationRepository) SaveCertification(ctx context.Context, certification domain.Certification) error {
	var certifications []*ent.CertificationCreate
	for _, hash := range certification.Hashes() {
		certifications = append(certifications, s.connection.DB().
			Certification.Create().
			SetHash(hash).
			SetAnchorID(certification.AnchorID()).
			SetAnchor(certification.Anchor()))
	}
	if _, err := s.connection.DB().Certification.CreateBulk(certifications...).Save(ctx); err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (s SQLCertificationRepository) GetCertification(ctx context.Context, anchor int, hash string) (*domain.Certification, error) {
	crt, err := s.connection.DB().Certification.Query().
		Where(certification.AnchorID(anchor), certification.And(certification.Hash(hash))).Only(ctx)
	if err != nil {
		return &domain.Certification{}, err
	}

	return domain.NewCertification(crt.AnchorID, []string{crt.Hash}, crt.Anchor), nil
}

func (s SQLCertificationRepository) UpdateCertificationAnchor(ctx context.Context, anchor integrity.Anchor) error {
	if _, err := s.connection.DB().Certification.Update().SetAnchor(&anchor).
		Where(certification.AnchorIDLTE(anchor.Id)).
		Save(ctx); err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}
