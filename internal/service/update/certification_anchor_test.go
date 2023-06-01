package update

import (
	"bloock-managed-api/internal/domain"
	mock_repository "bloock-managed-api/internal/domain/repository/mocks"
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
	notificationRepository := mock_repository.NewMockNotificationRepository(ctrl)
	anchor := &integrity.Anchor{
		Id:         int64(1),
		BlockRoots: []string{""},
		Networks:   []integrity.AnchorNetwork{},
		Root:       "root",
		Status:     "pending",
	}
	hash := "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"
	certification := domain.NewCertification(int(anchor.Id), hash, anchor)
	t.Run("given anchor it should be updated and send notification", func(t *testing.T) {
		certificationRepository.EXPECT().UpdateCertificationAnchor(context.TODO(), *anchor)
		certificationRepository.EXPECT().GetCertificationsByAnchorID(context.TODO(), int(anchor.Id)).
			Return([]domain.Certification{*certification}, nil)
		notificationRepository.EXPECT().NotifyCertification(hash, *anchor)

		err := NewCertificationAnchor(certificationRepository, notificationRepository).
			UpdateAnchor(context.TODO(), *anchor)

		assert.NoError(t, err)
	})

	t.Run("given error updating anchor it should be returned", func(t *testing.T) {
		certificationRepository.EXPECT().UpdateCertificationAnchor(context.TODO(), *anchor).
			Return(errors.New("some error"))

		err := NewCertificationAnchor(certificationRepository, notificationRepository).
			UpdateAnchor(context.TODO(), *anchor)

		assert.Error(t, err)
	})

	t.Run("given error notifying  it should be returned", func(t *testing.T) {
		certificationRepository.EXPECT().UpdateCertificationAnchor(context.TODO(), *anchor)
		certificationRepository.EXPECT().GetCertificationsByAnchorID(context.TODO(), int(anchor.Id)).
			Return([]domain.Certification{*certification}, nil)
		notificationRepository.EXPECT().NotifyCertification(hash, *anchor).Return(errors.New("some error"))

		err := NewCertificationAnchor(certificationRepository, notificationRepository).
			UpdateAnchor(context.TODO(), *anchor)

		assert.Error(t, err)
	})
}
