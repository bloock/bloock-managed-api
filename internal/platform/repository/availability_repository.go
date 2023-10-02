package repository

import (
	"context"
	"strings"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockAvailabilityRepository struct {
	recordClient       client.RecordClient
	availabilityClient client.AvailabilityClient
	logger             zerolog.Logger
}

func NewBloockAvailabilityRepository(logger zerolog.Logger) *BloockAvailabilityRepository {
	logger.With().Caller().Str("component", "availability-repository").Logger()
	return &BloockAvailabilityRepository{
		recordClient:       client.NewRecordClient(),
		availabilityClient: client.NewAvailabilityClient(),
		logger:             logger,
	}
}

func (b BloockAvailabilityRepository) UploadHosted(ctx context.Context, record *record.Record) (string, error) {
	return b.availabilityClient.Publish(*record, availability.NewHostedPublisher())
}

func (b BloockAvailabilityRepository) UploadIpfs(ctx context.Context, record *record.Record) (string, error) {
	return b.availabilityClient.Publish(*record, availability.NewIpfsPublisher())

}

func (b BloockAvailabilityRepository) FindFile(ctx context.Context, dataID string) ([]byte, error) {
	record, err := b.availabilityClient.Retrieve(availability.NewHostedLoader(dataID))
	if err != nil {
		record, err = b.availabilityClient.Retrieve(availability.NewIpfsLoader(dataID))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil, nil
			}
			return nil, err
		}
		return record.Retrieve(), nil
	}

	return record.Retrieve(), nil
}
