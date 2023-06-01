package repository

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBloockIntegrityRepository_Certify(t *testing.T) {
	data := [][]byte{
		[]byte("Hello World!"),
	}
	t.Run("given data to certify it should be certified with no errors", func(t *testing.T) {
		apiKey := "Nm1sFmrojcrRgfZ4v0H0w0d1d22GookjcJl7y-2jr51qx0RioCR3nVm1z74hDEzZ" //os.Getenv("BLOOCK_API_KEY")
		expectedHash := "3ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"
		data := data

		certification, err := NewBloockIntegrityRepository(apiKey, zerolog.Logger{}).Certify(context.TODO(), data)

		assert.NoError(t, err)
		assert.Equal(t, expectedHash, certification[0].Hash())
		assert.Greater(t, certification[0].AnchorID(), 0)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		apiKey := os.Getenv("MISS_API_KEY")
		data := data

		certification, err := NewBloockIntegrityRepository(apiKey, zerolog.Logger{}).Certify(context.TODO(), data)

		assert.Error(t, err)
		assert.Empty(t, certification)

	})
}
