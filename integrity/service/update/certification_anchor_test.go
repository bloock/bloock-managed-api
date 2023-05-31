package update

import (
	mock_repository "bloock-managed-api/integrity/domain/repository/mocks"
	"context"
	"errors"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCertificationAnchor_UpdateAnchor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	certificationRepository := mock_repository.NewMockCertificationRepository(ctrl)
	anchor := &integrity.Anchor{
		Id:         int64(1),
		BlockRoots: []string{""},
		Networks:   []integrity.AnchorNetwork{},
		Root:       "root",
		Status:     "pending",
	}
	t.Run("given anchor it should be updated", func(t *testing.T) {
		certificationRepository.EXPECT().UpdateCertificationAnchor(context.TODO(), *anchor)

		err := NewCertificationAnchor(certificationRepository).UpdateAnchor(context.TODO(), *anchor)

		assert.NoError(t, err)
	})

	t.Run("given error updating anchor it should be returned", func(t *testing.T) {
		certificationRepository.EXPECT().UpdateCertificationAnchor(context.TODO(), *anchor).Return(errors.New("some error"))

		err := NewCertificationAnchor(certificationRepository).UpdateAnchor(context.TODO(), *anchor)

		assert.Error(t, err)
	})
}
