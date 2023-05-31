package create

import (
	"bloock-managed-api/integrity/domain"
	mock_repository "bloock-managed-api/integrity/domain/repository/mocks"
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
		certification := domain.NewCertification(1, []string{
			"9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658",
			"4da432f1ecd4c0ac028ebde3a3f78510a21d54087b161590a63080d33b702b8d",
		}, nil)
		integrityRepository.EXPECT().Certify(ctx, files).Return(*certification, nil)
		certificationRepository.EXPECT().SaveCertification(ctx, *certification)

		err := NewCertification(certificationRepository, integrityRepository).Certify(ctx, files)

		assert.NoError(t, err)
	})

	t.Run("given error certifying it should be returned", func(t *testing.T) {
		certification := domain.NewCertification(1, []string{
			"9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658",
			"4da432f1ecd4c0ac028ebde3a3f78510a21d54087b161590a63080d33b702b8d",
		}, nil)
		integrityRepository.EXPECT().Certify(ctx, files).Return(*certification, nil)
		certificationRepository.EXPECT().SaveCertification(ctx, *certification).
			Return(someErr)

		err := NewCertification(certificationRepository, integrityRepository).Certify(ctx, files)

		assert.Error(t, err)
	})

	t.Run("given error saving certification result it should be returned", func(t *testing.T) {

		integrityRepository.EXPECT().Certify(ctx, files).Return(domain.Certification{}, someErr)

		err := NewCertification(certificationRepository, integrityRepository).Certify(ctx, files)

		assert.Error(t, err)
	})
}
