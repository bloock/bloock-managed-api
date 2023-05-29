package repository

import (
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockCertificationRepository struct {
	apikey          string
	integrityClient client.IntegrityClient
	log             zerolog.Logger
}

func NewBloockCertificationRepository(apikey string, log zerolog.Logger) *BloockCertificationRepository {
	bloock.ApiKey = apikey
	return &BloockCertificationRepository{
		apikey:          apikey,
		integrityClient: client.NewIntegrityClient(),
		log:             log,
	}
}

func (b BloockCertificationRepository) Certify(records []record.Record) (anchor int, err error) {

	receipt, err := b.integrityClient.SendRecords(records)
	if err != nil {
		b.log.Error().Err(err).Msg("error certifying data")
		return 0, err
	}

	return int(receipt[0].Anchor), nil
}
