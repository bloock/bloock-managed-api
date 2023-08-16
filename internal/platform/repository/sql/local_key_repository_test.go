package sql

import (
	"bloock-managed-api/internal/domain"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var ctx = context.TODO()

func TestSQLLocalKeyRepository_GetKeyByID(t *testing.T) {
	conn := EntConnectorForTesting(t)
	t.Run("given id localKey should be returned when exists", func(t *testing.T) {
		keyID := uuid.New()
		localKey := key.LocalKey{
			Key:        "k",
			PrivateKey: "p",
		}
		keyType := "Rsa2048"
		_, err := conn.DB().LocalKey.Create().
			SetID(keyID).
			SetLocalKey(&localKey).SetKeyType(keyType).
			Save(ctx)
		require.NoError(t, err)

		localKeyRepository := NewSQLLocalKeyRepository(*conn, 5*time.Second, zerolog.Logger{})
		keyByID, err := localKeyRepository.FindKeyByID(ctx, keyID)
		assert.NoError(t, err)
		assert.Equal(t, keyByID.Id(), keyID)
		assert.Equal(t, keyByID.LocalKey(), keyByID.LocalKey())
		assert.Equal(t, keyByID.KeyTypeStr(), keyType)
	})
}

func TestSQLLocalKeyRepository_SaveKey(t *testing.T) {
	conn := EntConnectorForTesting(t)
	localKey := key.LocalKey{
		Key:        "k",
		PrivateKey: "p",
	}
	t.Run("given localKey it should be saved correctly", func(t *testing.T) {
		localKeyRepository := NewSQLLocalKeyRepository(*conn, 5*time.Second, zerolog.Logger{})
		localKey := domain.NewLocalKey(localKey, key.Rsa4096, uuid.New())

		err := localKeyRepository.SaveKey(ctx, *localKey)

		require.NoError(t, err)
		currentLocalKey, err := conn.DB().LocalKey.Get(ctx, localKey.Id())
		require.NoError(t, err)
		assert.NotNil(t, currentLocalKey)
		assert.Equal(t, localKey.Id(), currentLocalKey.ID)
		assert.Equal(t, localKey.KeyTypeStr(), currentLocalKey.KeyType)
		assert.Equal(t, localKey.LocalKey(), *currentLocalKey.LocalKey)
	})
}

func TestSQLLocalKeyRepository_GetKeys(t *testing.T) {
	conn := EntConnectorForTesting(t)
	keyID := uuid.New()
	localKey := key.LocalKey{
		Key:        "k",
		PrivateKey: "p",
	}
	keyType := "Rsa2048"
	_, err := conn.DB().LocalKey.Create().
		SetID(keyID).
		SetLocalKey(&localKey).SetKeyType(keyType).
		Save(ctx)
	require.NoError(t, err)
	keyID = uuid.New()
	_, err = conn.DB().LocalKey.Create().
		SetID(keyID).
		SetLocalKey(&localKey).SetKeyType(keyType).
		Save(ctx)

	localKeyRepository := NewSQLLocalKeyRepository(*conn, 5*time.Second, zerolog.Logger{})
	localKeys, err := localKeyRepository.FindKeys(ctx)

	assert.NotEmpty(t, localKeys)

}
