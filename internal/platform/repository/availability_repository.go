package repository

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
)

type BloockAvailabilityRepository struct {
	recordClient       client.RecordClient
	availabilityClient client.AvailabilityClient
}

func NewBloockAvailabilityRepository(recordClient client.RecordClient, availabilityClient client.AvailabilityClient) *BloockAvailabilityRepository {
	return &BloockAvailabilityRepository{recordClient: recordClient, availabilityClient: availabilityClient}
}

func (b BloockAvailabilityRepository) UploadHosted(ctx context.Context, data []byte) (string, error) {
	rec, err := b.recordClient.FromBytes(data).Build()
	if err != nil {
		return "", err
	}
	return b.availabilityClient.Publish(rec, availability.NewHostedPublisher())
}

func (b BloockAvailabilityRepository) UploadIpfs(ctx context.Context, data []byte) (string, error) {
	rec, err := b.recordClient.FromBytes(data).Build()
	if err != nil {
		return "", err
	}
	return b.availabilityClient.Publish(rec, availability.NewIpfsPublisher())

}
