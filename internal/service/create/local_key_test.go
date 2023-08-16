package create

import (
	"bloock-managed-api/internal/domain"
	mock_repository "bloock-managed-api/internal/domain/repository/mocks"
	"errors"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLocalKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	localKeysRepository := mock_repository.NewMockLocalKeysRepository(ctrl)
	authenticityRepository := mock_repository.NewMockAuthenticityRepository(ctrl)
	t.Run("given key it should be created", func(t *testing.T) {
		keyType := key.EcP256k
		expectedLocalKey := domain.NewLocalKey(key.LocalKey{}, keyType, uuid.New())
		authenticityRepository.EXPECT().CreateLocalKey(keyType).Return(*expectedLocalKey, nil)
		localKeysRepository.EXPECT().SaveKey(gomock.Any(), *expectedLocalKey)
		localKey, err := NewLocalKey(authenticityRepository, localKeysRepository).Create(nil, keyType)

		assert.NoError(t, err)
		assert.NotNil(t, localKey)
	})

	t.Run("given key it should return error when creating key fails", func(t *testing.T) {
		keyType := key.EcP256k
		authenticityRepository.EXPECT().CreateLocalKey(keyType).Return(domain.LocalKey{}, errors.New("some err"))
		localKey, err := NewLocalKey(authenticityRepository, localKeysRepository).Create(nil, keyType)

		assert.Error(t, err)
		assert.Empty(t, localKey)
	})

	t.Run("given key it should return error when saving key fails", func(t *testing.T) {
		keyType := key.EcP256k
		expectedLocalKey := domain.NewLocalKey(key.LocalKey{}, keyType, uuid.New())
		authenticityRepository.EXPECT().CreateLocalKey(keyType).Return(*expectedLocalKey, nil)
		localKeysRepository.EXPECT().SaveKey(gomock.Any(), *expectedLocalKey).Return(errors.New("some err"))

		localKey, err := NewLocalKey(authenticityRepository, localKeysRepository).Create(nil, keyType)

		assert.Error(t, err)
		assert.Empty(t, localKey)
	})
}
