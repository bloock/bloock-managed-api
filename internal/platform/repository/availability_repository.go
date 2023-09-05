package repository

import (
	"context"
	"errors"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
	"github.com/rs/zerolog"
	"strings"
)

type BloockAvailabilityRepository struct {
	recordClient       client.RecordClient
	availabilityClient client.AvailabilityClient
	logger             zerolog.Logger
}

func NewBloockAvailabilityRepository(logger zerolog.Logger) *BloockAvailabilityRepository {
	logger.With().Caller().Str("component", "availability-repository").Logger()
	return &BloockAvailabilityRepository{
		recordClient: client.NewRecordClient(),
		availabilityClient: client.NewAvailabilityClient(),
		logger: logger,
	}
}

func (b BloockAvailabilityRepository) UploadHosted(ctx context.Context, data []byte) (string, error) {
	rec, err := b.recordClient.FromBytes(data).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", errUnknown
	}
	return b.availabilityClient.Publish(rec, availability.NewHostedPublisher())
}

func (b BloockAvailabilityRepository) UploadIpfs(ctx context.Context, data []byte) (string, error) {
	rec, err := b.recordClient.FromBytes(data).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", errUnknown
	}
	return b.availabilityClient.Publish(rec, availability.NewIpfsPublisher())

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

var errUnknown = errors.New("availability unknown error")
