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
	hash := "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"
	certification := domain.Certification {
		AnchorID: 1,
		Hash: hash,
		Data: file,
	}

	t.Run("given files it should be certified and saved in db", func(t *testing.T) {
		integrityRepository.EXPECT().Certify(gomock.Any(), file).Return(certification, nil)
		certificationRepository.EXPECT().SaveCertification(gomock.Any(), certification)

		_, err := NewIntegrityService(certificationRepository, integrityRepository).CertifyData(context.Background(), file)

		assert.NoError(t, err)
	})

	t.Run("given error saving certification result it should be returned", func(t *testing.T) {
		integrityRepository.EXPECT().Certify(gomock.Any(), file).Return(certification, nil)
		certificationRepository.EXPECT().SaveCertification(gomock.Any(), certification).
			Return(someErr)

		_, err := NewIntegrityService(certificationRepository, integrityRepository).CertifyData(context.Background(), file)

		assert.Error(t, err)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		integrityRepository.EXPECT().Certify(gomock.Any(), file).Return(domain.Certification{}, someErr)

		_, err := NewIntegrityService(certificationRepository, integrityRepository).CertifyData(context.Background(), file)

		assert.Error(t, err)
	})
}
