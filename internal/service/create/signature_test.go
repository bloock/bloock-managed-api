package create

import (
	"bloock-managed-api/internal/domain"
	mock_repository "bloock-managed-api/internal/domain/repository/mocks"
	"bloock-managed-api/internal/service/create/request"
	"errors"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authenticityRepository := mock_repository.NewMockAuthenticityRepository(ctrl)
	localKeysRepository := mock_repository.NewMockLocalKeysRepository(ctrl)
	t.Run("Given a sign request it should sign data and return signed data", func(t *testing.T) {
		commonName := "common name"
		keyType := key.Aes256
		dLocalKey := domain.NewLocalKey(key.LocalKey{
			Key:        "key",
			PrivateKey: "privKey",
		}, keyType, uuid.New())
		data := []byte("some_data")
		signRequest := request.NewSignRequest(dLocalKey.Id(), &commonName, data)
		expectedRecord, err := client.NewRecordClient().FromBytes(data).Build()
		require.NoError(t, err)
		localKeysRepository.EXPECT().FindKeyByID(gomock.Any(), dLocalKey.Id()).Return(dLocalKey, nil)
		authenticityRepository.EXPECT().Sign(dLocalKey.LocalKey(), dLocalKey.KeyType(), &commonName, data).Return(expectedRecord, nil)

		signResponse, err := NewSignature(authenticityRepository, localKeysRepository).Sign(nil, *signRequest)

		assert.NoError(t, err)
		assert.NotEmpty(t, signResponse.Record())
		assert.Equal(t, expectedRecord.Retrieve(), signResponse.Record())
	})

	t.Run("Given a sign request it should return error when sign fails", func(t *testing.T) {
		commonName := "common name"
		localKey := domain.NewLocalKey(key.LocalKey{
			Key:        "key",
			PrivateKey: "privKey",
		}, key.Aes256, uuid.New())
		data := []byte("some_data")
		signRequest := request.NewSignRequest(localKey.Id(), &commonName, data)
		localKeysRepository.EXPECT().FindKeyByID(gomock.Any(), localKey.Id()).Return(localKey, nil)
		authenticityRepository.EXPECT().Sign(localKey.LocalKey(), localKey.KeyType(), &commonName, data).Return(
			record.Record{},
			errors.New("some err"),
		)

		signResponse, err := NewSignature(authenticityRepository, localKeysRepository).Sign(nil, *signRequest)

		assert.Error(t, err)
		assert.Nil(t, signResponse)
	})

	t.Run("Given a sign request it should return error when getting keys fails", func(t *testing.T) {
		commonName := "common name"
		localKey := domain.NewLocalKey(key.LocalKey{
			Key:        "key",
			PrivateKey: "privKey",
		}, key.Aes256, uuid.New())
		data := []byte("some_data")
		signRequest := request.NewSignRequest(localKey.Id(), &commonName, data)
		localKeysRepository.EXPECT().FindKeyByID(gomock.Any(), localKey.Id()).Return(nil, errors.New("some err"))
		authenticityRepository.EXPECT().Sign(localKey.LocalKey(), localKey.KeyType(), &commonName, data).Times(0)

		signResponse, err := NewSignature(authenticityRepository, localKeysRepository).Sign(nil, *signRequest)

		assert.Error(t, err)
		assert.Nil(t, signResponse)
	})
}
