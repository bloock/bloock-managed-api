package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/pkg"
	"github.com/hashicorp/go-getter"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockAvailabilityRepository struct {
	client               client.BloockClient
	localStoragePath     string
	localStorageStrategy domain.LocalStorageStrategy
	tmpPath              string
	timeout              time.Duration
	logger               zerolog.Logger
}

func NewBloockAvailabilityRepository(ctx context.Context, l zerolog.Logger) repository.AvailabilityRepository {
	logger := l.With().Caller().Str("component", "availability-repository").Logger()

	c := client.NewBloockClient(pkg.GetApiKeyFromContext(ctx), nil)

	return &BloockAvailabilityRepository{
		client:               c,
		localStoragePath:     config.Configuration.Storage.LocalPath,
		localStorageStrategy: domain.LocalStorageStrategyFromString(config.Configuration.Storage.LocalStrategy),
		tmpPath:              config.Configuration.Storage.TmpDir,
		timeout:              15 * time.Second,
		logger:               logger,
	}
}

func (b BloockAvailabilityRepository) UploadHosted(ctx context.Context, file *domain.File, record record.Record) (string, error) {
	return b.client.AvailabilityClient.Publish(record, availability.NewHostedPublisher())
}

func (b BloockAvailabilityRepository) UploadIpfs(ctx context.Context, file *domain.File, record record.Record) (string, error) {
	return b.client.AvailabilityClient.Publish(record, availability.NewIpfsPublisher())
}

func (b BloockAvailabilityRepository) UploadLocal(ctx context.Context, file *domain.File) (string, error) {
	record, err := b.client.RecordClient.FromBytes(file.Bytes()).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", err
	}

	var name string
	switch b.localStorageStrategy {
	case domain.LocalStorageStrategyHash:
		hash, err := record.GetHash()
		if err != nil {
			if !errors.Is(err, os.ErrExist) {
				b.logger.Log().Err(err).Msg("")
				return "", errors.New("error retrieving record hash")
			}
		}

		name = hash
	case domain.LocalStorageStrategyFilename:
		name = fmt.Sprintf("%s%s", file.Filename(), file.FileExtension())
	default:
		return "", errors.New("invalid local storage strategy defined ")
	}

	return b.saveLocalFile(ctx, b.localStoragePath, name, record)
}

func (b BloockAvailabilityRepository) UploadTmp(ctx context.Context, file *domain.File, record record.Record) (string, error) {
	hash, err := record.GetHash()
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			b.logger.Log().Err(err).Msg("")
			return "", errors.New("error retrieving record hash")
		}
	}

	return b.saveLocalFile(ctx, b.tmpPath, hash, record)
}

func (b BloockAvailabilityRepository) RetrieveLocal(ctx context.Context, filePath string) ([]byte, error) {
	file, err := os.ReadFile(fmt.Sprintf("%s", filePath))
	if err != nil {
		return nil, errors.New("error retrieving the file")
	}

	return file, nil
}

func (b BloockAvailabilityRepository) RetrieveTmp(ctx context.Context, filename string) ([]byte, error) {

	file, err := os.ReadFile(fmt.Sprintf("%s/%s", b.tmpPath, filename))
	if err != nil {
		return nil, errors.New("error retrieving the file")
	}

	return file, nil
}

func (b BloockAvailabilityRepository) FindFile(ctx context.Context, id string) ([]byte, error) {
	if _, err := url.ParseRequestURI(id); err != nil {
		// is not a url
		file, err := b.downloadUrl(ctx, fmt.Sprintf("%s/hosting/v1/hosted/%s", config.Configuration.Bloock.CdnHost, id))
		if err != nil {
			file, err := b.downloadUrl(ctx, fmt.Sprintf("%s/hosting/v1/ipfs/%s", config.Configuration.Bloock.CdnHost, id))
			if err != nil {
				return nil, err
			}
			return file[0].Bytes(), nil
		}

		return file[0].Bytes(), nil
	} else {
		// is a url
		files, err := b.downloadUrl(ctx, id)
		if err != nil {
			return nil, err
		}
		return files[0].Bytes(), nil
	}
}

func (b BloockAvailabilityRepository) FindAll(ctx context.Context, id string) ([]domain.File, error) {
	return b.downloadUrl(ctx, id)
}

func (b BloockAvailabilityRepository) saveLocalFile(ctx context.Context, dir string, name string, record record.Record) (string, error) {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			b.logger.Log().Err(err).Msg("")
			return "", errors.New("error creating directory")
		}
	}

	fileBytes := record.Retrieve()
	path := fmt.Sprintf("%s/%s", dir, name)
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

func (b BloockAvailabilityRepository) downloadUrl(ctx context.Context, url string) ([]domain.File, error) {
	tmpDir := "./tmp_url"

	ctxTimeout, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	getterClient := &getter.Client{
		Ctx:  ctx,
		Dst:  tmpDir,
		Src:  url,
		Mode: getter.ClientModeAny,
		Getters: map[string]getter.Getter{
			"git":   &getter.GitGetter{},
			"https": &getter.HttpGetter{},
		},
		Options: []getter.ClientOption{
			getter.WithContext(ctxTimeout),
		},
	}
	if err := getterClient.Get(); err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("error downloading from url: %s", err.Error())
	}
	filesDomain := make([]domain.File, 0)

	if err := readFilesRecursive(tmpDir, &filesDomain); err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, err
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		return nil, err
	}

	return filesDomain, nil
}

func readFilesRecursive(dirPath string, filesDomain *[]domain.File) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		files = append(files, info)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		if file.IsDir() {
			err = readFilesRecursive(filePath, filesDomain)
			if err != nil {
				return err
			}
		} else {
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", filePath, err)
			}

			*filesDomain = append(*filesDomain, domain.NewFile(fileBytes, file.Name(), http.DetectContentType(fileBytes)))
		}
	}

	return nil
}
