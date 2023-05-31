package sql

import (
	"bloock-managed-api/ent/certification"
	"bloock-managed-api/integrity/domain"
	"bloock-managed-api/integrity/platform/repository/sql/connection"
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
	entConnector := connection.NewEntConnector(zerolog.Logger{})
	conn, err := connection.NewConnection("file:ent?mode=memory&cache=shared&_fk=1", entConnector, zerolog.Logger{})
	require.NoError(t, err)
	err = conn.Migrate()
	require.NoError(t, err)

	certificationRepository := NewSQLCertificationRepository(*conn, 5*time.Second, zerolog.Logger{})
	t.Run("given certification it should be saved", func(t *testing.T) {
		certification := domain.NewCertification(1, []string{hash}, nil)
		err := certificationRepository.SaveCertification(context.TODO(), *certification)

		assert.NoError(t, err)
	})

}

func TestSQLCertificationRepository_GetCertification(t *testing.T) {
	entConnector := connection.NewEntConnector(zerolog.Logger{})
	conn, err := connection.NewConnection("file:ent?mode=memory&cache=shared&_fk=1", entConnector, zerolog.Logger{})
	require.NoError(t, err)
	err = conn.Migrate()
	require.NoError(t, err)

	certificationRepository := NewSQLCertificationRepository(*conn, 5*time.Second, zerolog.Logger{})
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

		certification, err := certificationRepository.GetCertification(ctx, anchorID, hash)

		assert.NoError(t, err)
		assert.Equal(t, hash, certification.Hashes()[0])
		assert.Equal(t, anchorID, certification.AnchorID())
		assert.Equal(t, anchor, certification.Anchor())
	})

	t.Run("given anchor and hash it should return error when doesnt exists", func(t *testing.T) {

		certification, err := certificationRepository.GetCertification(ctx, 2, hash)

		assert.Error(t, err)
		assert.Empty(t, certification)
	})
}

func TestSQLCertificationRepository_UpdateCertificationAnchor(t *testing.T) {
	entConnector := connection.NewEntConnector(zerolog.Logger{})
	conn, err := connection.NewConnection("file:ent?mode=memory&cache=shared&_fk=1", entConnector, zerolog.Logger{})
	require.NoError(t, err)
	err = conn.Migrate()
	require.NoError(t, err)
	hash := "2ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"

	certificationRepository := NewSQLCertificationRepository(*conn, 5*time.Second, zerolog.Logger{})
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
