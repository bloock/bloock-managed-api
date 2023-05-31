package retrieve

import (
	"bloock-managed-api/integrity/domain"
	mock_repository "bloock-managed-api/integrity/domain/repository/mocks"
	"context"
	"errors"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCertificationByAnchorAndHash_Retrieve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	certificationRepository := mock_repository.NewMockCertificationRepository(ctrl)

	t.Run("given anchorID and hash it should retrieve a certification", func(t *testing.T) {
		hash := "3ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"
		anchorID := 1
		anchor := &integrity.Anchor{}
		certification := domain.NewCertification(anchorID, []string{hash}, anchor)
		expectedResponse := NewCertificationByAnchorAndHashResponse(hash, anchor)
		certificationRepository.EXPECT().
			GetCertification(context.TODO(), anchorID, hash).
			Return(certification, nil)

		certificationByAnchorAndHashResponse, err := NewCertificationByAnchorAndHash(certificationRepository).Retrieve(context.TODO(), anchorID, hash)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, certificationByAnchorAndHashResponse)
	})

	t.Run("given an error it should be returned", func(t *testing.T) {
		hash := "3ea2f1d0abf3fc66cf29eebb70cbd4e7fe762ef8a09bcc06c8edf641230afec0"
		anchorID := 1
		certificationRepository.EXPECT().
			GetCertification(context.TODO(), anchorID, hash).
			Return(&domain.Certification{}, errors.New("some error"))

		_, err := NewCertificationByAnchorAndHash(certificationRepository).
			Retrieve(context.TODO(), anchorID, hash)

		assert.Error(t, err)

	})
}
