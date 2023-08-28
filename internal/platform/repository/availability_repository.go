package repository

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
	"github.com/rs/zerolog"
)

type BloockAvailabilityRepository struct {
	recordClient       client.RecordClient
	availabilityClient client.AvailabilityClient
	logger             zerolog.Logger
}

func NewBloockAvailabilityRepository(recordClient client.RecordClient, availabilityClient client.AvailabilityClient, logger zerolog.Logger) *BloockAvailabilityRepository {
	return &BloockAvailabilityRepository{recordClient: recordClient, availabilityClient: availabilityClient, logger: logger}
}

func (b BloockAvailabilityRepository) UploadHosted(ctx context.Context, data []byte) (string, error) {
	rec, err := b.recordClient.FromBytes(data).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", err
	}
	return b.availabilityClient.Publish(rec, availability.NewHostedPublisher())
}

func (b BloockAvailabilityRepository) UploadIpfs(ctx context.Context, data []byte) (string, error) {
	rec, err := b.recordClient.FromBytes(data).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", err
	}
	return b.availabilityClient.Publish(rec, availability.NewIpfsPublisher())

}
