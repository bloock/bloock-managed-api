package repository

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/availability"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var apiKey = "Nm1sFmrojcrRgfZ4v0H0w0d1d22GookjcJl7y-2jr51qx0RioCR3nVm1z74hDEzZ"

func TestBloockAvailabilityRepository_UploadHosted(t *testing.T) {
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

func TestBloockAvailabilityRepository_FindFile(t *testing.T) {
	bloock.ApiKey = apiKey

	t.Run("given data id it should return data if it exists", func(t *testing.T) {
		record, err := client.NewRecordClient().FromString("Hello").Build()
		require.NoError(t, err)
		availabilityClient := client.NewAvailabilityClient()
		expectedDataID, err := availabilityClient.Publish(record, availability.NewHostedPublisher())
		require.NoError(t, err)

		currentDataBytes, err := NewBloockAvailabilityRepository(
			client.NewRecordClient(),
			client.NewAvailabilityClient(),
			zerolog.Logger{},
		).FindFile(context.TODO(), expectedDataID)

		assert.Equal(t, record.Retrieve(), currentDataBytes)
	})

	t.Run("given dataID when data is not found error shouldn't be returned", func(t *testing.T) {
		currentDataBytes, err := NewBloockAvailabilityRepository(
			client.NewRecordClient(),
			client.NewAvailabilityClient(),
			zerolog.Logger{},
		).FindFile(context.TODO(), "f7df610d-9f5e-446a-bc4f-1d3012410490")

		assert.NoError(t, err)
		assert.Empty(t, currentDataBytes)
	})

	t.Run("given dataID when error occurs it should be returned", func(t *testing.T) {
		currentDataBytes, err := NewBloockAvailabilityRepository(
			client.NewRecordClient(),
			client.NewAvailabilityClient(),
			zerolog.Logger{},
		).FindFile(context.TODO(), "")

		assert.Error(t, err)
		assert.ErrorIs(t, err, errUnknown)
		assert.Empty(t, currentDataBytes)
	})

}
