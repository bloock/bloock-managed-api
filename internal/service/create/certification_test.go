package create

import (
	"bloock-managed-api/internal/domain"
	mock_repository "bloock-managed-api/internal/domain/repository/mocks"
	"errors"
	"github.com/stretchr/testify/assert"

	"context"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCertification_Certify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	certificationRepository := mock_repository.NewMockCertificationRepository(ctrl)
	integrityRepository := mock_repository.NewMockIntegrityRepository(ctrl)
	someErr := errors.New("some error")

	files := [][]byte{[]byte("test"), []byte("test2")}
	ctx := context.TODO()
	t.Run("given files it should be certified and saved in db", func(t *testing.T) {
		certification := domain.NewPendingCertification(1, "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658")
		integrityRepository.EXPECT().Certify(ctx, files).Return([]domain.Certification{*certification}, nil)
		certificationRepository.EXPECT().SaveCertification(ctx, []domain.Certification{*certification})

		_, err := NewCertification(certificationRepository, integrityRepository).Certify(ctx, files)

		assert.NoError(t, err)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		certification := domain.NewPendingCertification(1, "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658")
		integrityRepository.EXPECT().Certify(ctx, files).Return([]domain.Certification{*certification}, nil)
		certificationRepository.EXPECT().SaveCertification(ctx, []domain.Certification{*certification}).
			Return(someErr)

		_, err := NewCertification(certificationRepository, integrityRepository).Certify(ctx, files)

		assert.Error(t, err)
	})

	t.Run("given error saving certification result it should be returned", func(t *testing.T) {

		integrityRepository.EXPECT().Certify(ctx, files).Return([]domain.Certification{}, someErr)

		_, err := NewCertification(certificationRepository, integrityRepository).Certify(ctx, files)

		assert.Error(t, err)
	})
}
