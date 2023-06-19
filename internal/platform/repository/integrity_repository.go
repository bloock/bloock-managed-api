package repository

import (
	"bloock-managed-api/internal/domain"
	"context"
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockIntegrityRepository struct {
	apikey          string
	integrityClient client.IntegrityClient
	log             zerolog.Logger
}

func NewBloockIntegrityRepository(apikey string, log zerolog.Logger) *BloockIntegrityRepository {
	bloock.ApiKey = apikey
	return &BloockIntegrityRepository{
		apikey:          apikey,
		integrityClient: client.NewIntegrityClient(),
		log:             log,
	}
}

func (b BloockIntegrityRepository) Certify(ctx context.Context, files [][]byte) (certification []domain.Certification, err error) {
	var records []record.Record
	for i := range files {
		rec, err := client.NewRecordClient().FromBytes(files[i]).Build()
		if err != nil {
			b.log.Error().Err(err).Msg("error certifying data")
			return []domain.Certification{}, err
		}

		records = append(records, rec)
	}

	receipt, err := b.integrityClient.SendRecords(records)
	if err != nil {
		b.log.Error().Err(err).Msg("error certifying data")
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
		b.log.Error().Err(err).Msg("error getting anchor")
		return integrity.Anchor{}, err
	}

	return anchor, nil
}
