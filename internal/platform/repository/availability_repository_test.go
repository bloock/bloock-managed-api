package repository

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBloockAvailabilityRepository_UploadHosted(t *testing.T) {
	apiKey := "Nm1sFmrojcrRgfZ4v0H0w0d1d22GookjcJl7y-2jr51qx0RioCR3nVm1z74hDEzZ"
	bloock.ApiKey = apiKey

	t.Run("given data it should be uploaded", func(t *testing.T) {
		id, err := NewBloockAvailabilityRepository(client.NewRecordClient(), client.NewAvailabilityClient(), zerolog.Logger{}).UploadHosted(context.TODO(), []byte("Hello World"))

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("given error it should be returned", func(t *testing.T) {
		bloock.ApiKey = ""
		_, err := NewBloockAvailabilityRepository(client.NewRecordClient(), client.NewAvailabilityClient(), zerolog.Logger{}).UploadHosted(context.TODO(), []byte("Hello World"))

		assert.Error(t, err)
	})
}

func TestBloockAvailabilityRepository_UploadIpfs(t *testing.T) {
	apiKey := "Nm1sFmrojcrRgfZ4v0H0w0d1d22GookjcJl7y-2jr51qx0RioCR3nVm1z74hDEzZ"
	bloock.ApiKey = apiKey

	t.Run("given data it should be uploaded", func(t *testing.T) {
		id, err := NewBloockAvailabilityRepository(client.NewRecordClient(), client.NewAvailabilityClient(), zerolog.Logger{}).UploadIpfs(context.TODO(), []byte("Hello World"))

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("given error it should be returned", func(t *testing.T) {
		bloock.ApiKey = ""
		_, err := NewBloockAvailabilityRepository(client.NewRecordClient(), client.NewAvailabilityClient(), zerolog.Logger{}).UploadIpfs(context.TODO(), []byte("Hello World"))

		assert.Error(t, err)
	})
}
