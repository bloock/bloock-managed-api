package hard_drive

import (
	"bloock-managed-api/internal/platform/test_utils/fixtures"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestHardDriveLocalStorageRepository_Save(t *testing.T) {
	t.Run("given data and filename it should be saved in hard drive", func(t *testing.T) {
		directory := "tmp"
		hash := "testFile"

		err := NewHardDriveLocalStorageRepository(directory, zerolog.Logger{}).
			Save(context.TODO(), fixtures.PDFContent, hash)

		fullPath := fmt.Sprintf("%s/%s", directory, hash)
		file, err := os.ReadFile(fullPath)
		require.NoError(t, err)

		assert.NoError(t, err)
		assert.NotEmpty(t, file)

		err = os.Remove(fullPath)
		require.NoError(t, err)
		err = os.Remove(directory)
		require.NoError(t, err)
	})

}

func TestHardDriveLocalStorageRepository_Retrieve(t *testing.T) {
	directory := "../../test_utils/fixtures/"
	filename := "pdf.pdf"

	data, err := NewHardDriveLocalStorageRepository(directory, zerolog.Logger{}).
		Retrieve(context.TODO(), directory, filename)

	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Equal(t, fixtures.PDFContent, data)

}
