package repository

import (
	"bloock-managed-api/integrity/domain"
	"context"
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
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
	for i, _ := range files {
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
