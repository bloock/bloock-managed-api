package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockAvailabilityRepository struct {
	recordClient       client.RecordClient
	availabilityClient client.AvailabilityClient
	localStoragePath   string
	tmpPath            string
	logger             zerolog.Logger
}

func NewBloockAvailabilityRepository(localStoragePath, tmpPath string, logger zerolog.Logger) *BloockAvailabilityRepository {
	logger.With().Caller().Str("component", "availability-repository").Logger()
	return &BloockAvailabilityRepository{
		recordClient:       client.NewRecordClient(),
		availabilityClient: client.NewAvailabilityClient(),
		localStoragePath:   localStoragePath,
		tmpPath:            tmpPath,
		logger:             logger,
	}
}

func (b BloockAvailabilityRepository) UploadHosted(ctx context.Context, record *record.Record) (string, error) {
	return b.availabilityClient.Publish(*record, availability.NewHostedPublisher())
}

func (b BloockAvailabilityRepository) UploadIpfs(ctx context.Context, record *record.Record) (string, error) {
	return b.availabilityClient.Publish(*record, availability.NewIpfsPublisher())
}

func (b BloockAvailabilityRepository) UploadLocal(ctx context.Context, record *record.Record) (string, error) {
	hash, err := record.GetHash()
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			b.logger.Log().Err(err).Msg("")
			return "", errors.New("error retrieving record hash")
		}
	}

	return b.saveLocalFile(ctx, b.localStoragePath, hash, record)
}

func (b BloockAvailabilityRepository) UploadTmp(ctx context.Context, record *record.Record) (string, error) {
	hash, err := record.GetHash()
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			b.logger.Log().Err(err).Msg("")
			return "", errors.New("error retrieving record hash")
		}
	}

	return b.saveLocalFile(ctx, b.tmpPath, hash, record)
}

func (b BloockAvailabilityRepository) RetrieveTmp(ctx context.Context, filename string) ([]byte, error) {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", b.tmpPath, filename))
	if err != nil {
		b.logger.Log().Err(err).Msg("")
		return nil, errors.New("error retrieving the file")
	}

	return file, nil
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

func (b BloockAvailabilityRepository) saveLocalFile(ctx context.Context, dir string, name string, record *record.Record) (string, error) {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			b.logger.Log().Err(err).Msg("")
			return "", errors.New("error creating directory")
		}
	}

	hash, err := record.GetHash()
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			b.logger.Log().Err(err).Msg("")
			return "", errors.New("error retrieving record hash")
		}
	}

	fileBytes := record.Retrieve()
	path := fmt.Sprintf("%s/%s", dir, hash)
	if err = os.WriteFile(path, fileBytes, 0644); err != nil {
		b.logger.Log().Err(err).Msg("")
		return "", err
	}

	uri, err := url.Parse(path)
	if err != nil {
		b.logger.Log().Err(err).Msg("")
		return "", err
	}

	return uri.String(), nil
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
