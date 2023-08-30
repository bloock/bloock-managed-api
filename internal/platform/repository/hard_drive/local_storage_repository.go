package hard_drive

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

type HardDriveLocalStorageRepository struct {
	directory string
	logger    zerolog.Logger
}

func NewHardDriveLocalStorageRepository(directory string, logger zerolog.Logger) *HardDriveLocalStorageRepository {
	logger.With().Caller().Str("component", "hard-drive-local-storage").Logger()
	return &HardDriveLocalStorageRepository{directory: directory, logger: logger}
}

func (h HardDriveLocalStorageRepository) Save(ctx context.Context, data []byte, hash string) error {
	err := os.Mkdir(h.directory, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			h.logger.Log().Err(err).Msg("")
			msg := "error creating directory"
			err := errors.New(msg)
			return err
		}
	}

	if err := os.WriteFile(fmt.Sprintf("%s/%s", h.directory, hash), data, 0644); err != nil {
		h.logger.Log().Err(err).Msg("")
		return err
	}
	return nil
}

func (h HardDriveLocalStorageRepository) Retrieve(ctx context.Context, directory string, filename string) ([]byte, error) {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		h.logger.Log().Err(err).Msg("")
		msg := "error retrieving the file"
		err := errors.New(msg)
		return nil, err
	}

	return file, nil
}
