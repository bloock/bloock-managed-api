package repository

import (
	"bloock-managed-api/internal/platform/test_utils/fixtures"
	"context"
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBloockIntegrityRepository_Certify(t *testing.T) {
	data := fixtures.PDFContent

	apiKey := "Nm1sFmrojcrRgfZ4v0H0w0d1d22GookjcJl7y-2jr51qx0RioCR3nVm1z74hDEzZ"
	bloock.ApiKey = apiKey
	t.Run("given data to certify it should be certified with no errors", func(t *testing.T) {
		expectedHash := "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5"
		data := data

		certification, err := NewBloockIntegrityRepository(client.NewIntegrityClient(), zerolog.Logger{}).Certify(context.TODO(), data)

		assert.NoError(t, err)
		assert.Equal(t, expectedHash, certification.Hash())
		assert.Greater(t, certification.AnchorID(), 0)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		bloock.ApiKey = ""
		data := data

		certification, err := NewBloockIntegrityRepository(client.NewIntegrityClient(), zerolog.Logger{}).Certify(context.TODO(), data)

		assert.Error(t, err)
		assert.Empty(t, certification)

	})
}
