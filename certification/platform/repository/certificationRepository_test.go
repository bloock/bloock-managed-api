package repository

import (
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBloockCertificationRepository_Certify(t *testing.T) {
	t.Run("given data to certify it should be certified with no errors", func(t *testing.T) {
		apiKey := os.Getenv("BLOOCK_API_KEY")
		rec, err := client.NewRecordClient().FromBytes([]byte("Hello World!")).Build()
		require.NoError(t, err)
		data := []record.Record{rec}

		anchor, err := NewBloockCertificationRepository(apiKey, zerolog.Logger{}).Certify(data)

		assert.NoError(t, err)
		assert.Greater(t, anchor, 0)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		apiKey := os.Getenv("MISS_API_KEY")
		rec, err := client.NewRecordClient().FromBytes([]byte("Hello World!")).Build()
		require.NoError(t, err)
		data := []record.Record{rec}

		anchor, err := NewBloockCertificationRepository(apiKey, zerolog.Logger{}).Certify(data)

		assert.Error(t, err)
		assert.Equal(t, 0, anchor)
	})
}
