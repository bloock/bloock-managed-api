package availability

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
)

type AvailabilityService struct {
	availabilityClient client.AvailabilityClient
	recordClient       client.RecordClient
}

func (a AvailabilityService) UploadHosted(ctx context.Context, data []byte) (string, error) {
	rec, err := a.recordClient.FromFile(data).Build()
	if err != nil {
		return "", err
	}
	return a.availabilityClient.Publish(rec, availability.NewHostedPublisher())
}

func (a AvailabilityService) UploadIpfs(ctx context.Context, data []byte) (string, error) {
	rec, err := a.recordClient.FromFile(data).Build()
	if err != nil {
		return "", err
	}
	return a.availabilityClient.Publish(rec, availability.NewIpfsPublisher())
}
