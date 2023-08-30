package integrity

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

	file := []byte("test")
	ctx := context.TODO()
	t.Run("given files it should be certified and saved in db", func(t *testing.T) {
		hash := "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"
		certification := domain.NewPendingCertification(1, hash, file)
		integrityRepository.EXPECT().Certify(ctx, file).Return(*certification, nil)
		certificationRepository.EXPECT().SaveCertification(ctx, *certification)

		_, err := NewIntegrityService(certificationRepository, integrityRepository).Certify(ctx, file)

		assert.NoError(t, err)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		certification := domain.NewPendingCertification(1, "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658", file)
		integrityRepository.EXPECT().Certify(gomock.Any(), file).Return(*certification, nil)
		certificationRepository.EXPECT().SaveCertification(ctx, *certification).
			Return(someErr)

		_, err := NewIntegrityService(certificationRepository, integrityRepository).Certify(ctx, file)

		assert.Error(t, err)
	})

	t.Run("given error saving certification result it should be returned", func(t *testing.T) {

		integrityRepository.EXPECT().Certify(ctx, file).Return(domain.Certification{}, someErr)

		_, err := NewIntegrityService(certificationRepository, integrityRepository).Certify(ctx, file)

		assert.Error(t, err)
	})

	t.Run("given data id and hash certification should be updated", func(t *testing.T) {
		hash := "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"
		dataID := "fbde810d-0379-42b3-ac9a-eaba3c48fe1a"
		certificationRepository.EXPECT().
			UpdateCertificationDataID(
				ctx,
				hash,
				dataID,
			)

		err := NewIntegrityService(certificationRepository, integrityRepository).SetDataIDToCertification(ctx, hash, dataID)

		assert.NoError(t, err)
	})
}
