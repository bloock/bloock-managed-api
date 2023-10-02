package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

func (b BloockAvailabilityRepository) FindFile(ctx context.Context, id string) ([]byte, error) {
	if _, err := url.ParseRequestURI(id); err != nil {
		// is not a url

		file, err := b.downloadUrl(ctx, fmt.Sprintf("https://cdn.bloock.com/hosting/v1/hosted/%s", id))
		if err != nil {
			file, err := b.downloadUrl(ctx, fmt.Sprintf("https://cdn.bloock.com/hosting/v1/ipfs/%s", id))
			if err != nil {
				return nil, err
			}
			return file, nil
		}

		return file, nil
	} else {
		// is a url
		return b.downloadUrl(ctx, id)
	}
}

func (b BloockAvailabilityRepository) downloadUrl(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("error downloading file from %s: %s", url, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("error downloading file from %s: received status code %d", url, resp.StatusCode)
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error downloading file from %s: %s", url, err.Error())
	}

	return file, nil
}
