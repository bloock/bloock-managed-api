package authenticity

import (
	mock_repository "bloock-managed-api/internal/domain/repository/mocks"
	"bloock-managed-api/internal/service/authenticity/request"
	"errors"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateManagedKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticityRepository := mock_repository.NewMockAuthenticityRepository(ctrl)
	t.Run("given key it should be created", func(t *testing.T) {
		keyType := key.EcP256k
		protectionLevel := key.KEY_PROTECTION_SOFTWARE
		name := "test"
		expiration := 312300
		expectedLocalKey := key.ManagedKey{
			ID:         uuid.NewString(),
			Name:       "name",
			Protection: 1,
			KeyType:    1,
			Expiration: 10,
			Key:        "key",
		}
		authenticityRepository.EXPECT().CreateManagedKey(name, keyType, expiration, protectionLevel).Return(expectedLocalKey, nil)
		localKey, err := NewManagedKey(authenticityRepository).Create(*request.NewCreateManagedKeyRequest(name, keyType, expiration, protectionLevel))

		assert.NoError(t, err)
		assert.NotNil(t, localKey)
	})

	t.Run("given key it should return error when creating key fails", func(t *testing.T) {
		keyType := key.EcP256k
		protectionLevel := key.KEY_PROTECTION_SOFTWARE
		name := "test"
		expiration := 312300
		authenticityRepository.EXPECT().CreateManagedKey(name, keyType, expiration, protectionLevel).Return(key.ManagedKey{}, errors.New("some err"))

		localKey, err := NewManagedKey(authenticityRepository).Create(*request.NewCreateManagedKeyRequest(name, keyType, expiration, protectionLevel))

		assert.Error(t, err)
		assert.Empty(t, localKey)
	})

}
