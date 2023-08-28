package repository

import (
	"bloock-managed-api/internal/domain"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockIntegrityRepository struct {
	integrityClient client.IntegrityClient
	logger          zerolog.Logger
}

func NewBloockIntegrityRepository(integrityClient client.IntegrityClient, logger zerolog.Logger) *BloockIntegrityRepository {
	logger.With().Caller().Str("component", "integrity-repository").Logger()
	return &BloockIntegrityRepository{integrityClient: integrityClient, logger: logger}
}

func (b BloockIntegrityRepository) Certify(ctx context.Context, file []byte) (certification []domain.Certification, err error) {

	rec, err := client.NewRecordClient().FromFile(file).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("error certifying data")
		return []domain.Certification{}, err
	}

	receipt, err := b.integrityClient.SendRecords([]record.Record{rec})
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return []domain.Certification{}, err
	}

	var certifications []domain.Certification
	for _, recordReceipt := range receipt {
		crts := *domain.NewPendingCertification(int(recordReceipt.Anchor), recordReceipt.Record)
		certifications = append(certifications, crts)
	}

	return certifications, nil
}

func (b BloockIntegrityRepository) GetAnchorByID(ctx context.Context, anchorID int) (integrity.Anchor, error) {
	anchor, err := b.integrityClient.GetAnchor(int64(anchorID))
	if err != nil {
		b.logger.Error().Err(err).Msg("error getting anchor")
		return integrity.Anchor{}, err
	}

	return anchor, nil
}
