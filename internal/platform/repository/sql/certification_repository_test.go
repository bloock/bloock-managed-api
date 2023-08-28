package sql_test

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/repository/sql"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/repository/sql/ent/certification"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
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
		certification := domain.NewPendingCertification(1, hash)
		err := certificationRepository.SaveCertification(context.TODO(), []domain.Certification{*certification})

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
		anchor := &integrity.Anchor{
			Id:         int64(anchorID),
			BlockRoots: []string{""},
			Networks:   []integrity.AnchorNetwork{},
			Root:       "root",
			Status:     "pending",
		}
		_, err := conn.DB().Certification.Create().
			SetAnchor(anchor).
			SetHash(hash).
			SetAnchorID(anchorID).
			Save(ctx)
		require.NoError(t, err)

		certifications, err := certificationRepository.GetCertificationsByAnchorID(ctx, anchorID)

		assert.NoError(t, err)
		assert.Equal(t, hash, certifications[0].Hash())
		assert.Equal(t, anchorID, certifications[0].AnchorID())
		assert.Equal(t, anchor, certifications[0].Anchor())
	})

	t.Run("given anchor and hash it should return empty list when no certification exists", func(t *testing.T) {

		certifications, err := certificationRepository.GetCertificationsByAnchorID(ctx, 2)

		assert.NoError(t, err)
		assert.Empty(t, certifications)
	})
}

func TestSQLCertificationRepository_UpdateCertificationAnchor(t *testing.T) {
	conn := EntConnectorForTesting(t)
	hash := "2ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"

	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, zerolog.Logger{})
	anchorID := 3
	smallAnchorID := 2
	anchor := &integrity.Anchor{
		Id:         int64(anchorID),
		BlockRoots: []string{""},
		Networks:   []integrity.AnchorNetwork{},
		Root:       "root",
		Status:     "pending",
	}
	updatedAnchor := &integrity.Anchor{
		Id:         int64(anchorID),
		BlockRoots: []string{""},
		Networks:   []integrity.AnchorNetwork{},
		Root:       "root",
		Status:     "success",
	}

	t.Run("given anchor it should update certification anchor for all smallest anchorIds", func(t *testing.T) {
		_, err := conn.DB().Certification.Create().
			SetAnchor(anchor).
			SetHash(hash).
			SetAnchorID(anchorID).
			Save(context.TODO())
		require.NoError(t, err)
		_, err = conn.DB().Certification.Create().
			SetAnchor(anchor).
			SetHash(hash).
			SetAnchorID(smallAnchorID).
			Save(context.TODO())
		require.NoError(t, err)

		err = certificationRepository.UpdateCertificationAnchor(context.TODO(), *updatedAnchor)

		certifications, queryErr := conn.DB().Certification.Query().
			Where(certification.Hash(hash)).All(context.TODO())
		require.NoError(t, queryErr)
		assert.NoError(t, err)
		assert.Equal(t, updatedAnchor, certifications[0].Anchor)
		assert.Equal(t, anchorID, certifications[0].AnchorID)
		assert.Equal(t, hash, certifications[0].Hash)
		assert.Equal(t, updatedAnchor, certifications[1].Anchor)
	})
}
