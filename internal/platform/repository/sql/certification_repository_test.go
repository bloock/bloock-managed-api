package sql_test

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/repository/sql"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"time"

	"testing"
)

var hash = "3ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"

func TestSQLCertificationRepository_SaveCertification(t *testing.T) {
	conn := EntConnectorForTesting(t)

	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, zerolog.Logger{})
	t.Run("given certification it should be saved", func(t *testing.T) {
		certification := domain.Certification{
			AnchorID: 1,
			Hash: hash,
			Data: nil,
		}
		err := certificationRepository.SaveCertification(context.Background(), certification)

		assert.NoError(t, err)
	})

}

func EntConnectorForTesting(t *testing.T) *connection.EntConnection {
	entConnector := connection.NewEntConnector(zerolog.Logger{})
	conn, err := connection.NewEntConnection("file:ent?mode=memory&cache=shared&_fk=1", entConnector, zerolog.Logger{})
	require.NoError(t, err)
	err = conn.Migrate()
	require.NoError(t, err)
	return conn
}

func TestSQLCertificationRepository_GetCertification(t *testing.T) {
	conn := EntConnectorForTesting(t)
	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, zerolog.Logger{})
	hash := "3ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"
	ctx := context.TODO()
	t.Run("given anchor and hash it should be returned when exists", func(t *testing.T) {
		anchorID := 5
		_, err := conn.DB().Certification.Create().
			SetHash(hash).
			SetAnchorID(anchorID).
			SetDataID("").
			Save(ctx)
		require.NoError(t, err)

		certifications, err := certificationRepository.GetCertificationsByAnchorID(ctx, anchorID)

		assert.NoError(t, err)
		assert.Equal(t, hash, certifications[0].Hash)
		assert.Equal(t, anchorID, certifications[0].AnchorID)
	})

	t.Run("given anchor and hash it should return empty list when no certification exists", func(t *testing.T) {
		certifications, err := certificationRepository.GetCertificationsByAnchorID(ctx, 2)

		assert.NoError(t, err)
		assert.Empty(t, certifications)
	})
}

func TestSQLCertificationRepository_UpdateCertificationDataID(t *testing.T) {
	conn := EntConnectorForTesting(t)
	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, zerolog.Logger{})

	t.Run("given hash and data id it should update data id for existent certification", func(t *testing.T) {
		crt, err := conn.DB().Certification.Create().
			SetHash(hash).
			SetAnchorID(1).
			SetDataID("").
			Save(context.TODO())
		require.NoError(t, err)

		updateCertification := domain.Certification{
			AnchorID: 1,
			Hash: hash,
			DataID: "a47ef5f4-26ba-4bbe-b53a-2e73a4d69001",
		}
		err = certificationRepository.UpdateCertificationDataID(ctx, updateCertification)
		assert.NoError(t, err)

		expectedCrt, err := conn.DB().Certification.Get(ctx, crt.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, expectedCrt.DataID)
		assert.Equal(t, expectedCrt.DataID, updateCertification.DataID)

		exists, err := certificationRepository.ExistCertificationByHash(context.Background(), updateCertification.Hash)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}
