package repository

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBloockIntegrityRepository_Certify(t *testing.T) {
	data := []byte("Hello World!")

	apiKey := "Nm1sFmrojcrRgfZ4v0H0w0d1d22GookjcJl7y-2jr51qx0RioCR3nVm1z74hDEzZ"
	bloock.ApiKey = apiKey
	t.Run("given data to certify it should be certified with no errors", func(t *testing.T) {
		expectedHash := "3ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"
		data := data

		certification, err := NewBloockIntegrityRepository(client.NewIntegrityClient(), zerolog.Logger{}).Certify(context.TODO(), data)

		assert.NoError(t, err)
		assert.Equal(t, expectedHash, certification[0].Hash())
		assert.Greater(t, certification[0].AnchorID(), 0)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		bloock.ApiKey = ""
		data := data

		certification, err := NewBloockIntegrityRepository(client.NewIntegrityClient(), zerolog.Logger{}).Certify(context.TODO(), data)

		assert.Error(t, err)
		assert.Empty(t, certification)

	})
}
