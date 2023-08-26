package availability_test

import (
	mock_repository "bloock-managed-api/internal/domain/repository/mocks"
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/availability"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAvailabilityService_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	availabilityRepository := mock_repository.NewMockAvailabilityRepository(ctrl)

	t.Run("given data it should be published", func(t *testing.T) {
		data := []byte("hello")
		response := "url"
		availabilityRepository.EXPECT().UploadHosted(gomock.Any(), data).Return(response, nil)

		url, err := availability.NewAvailabilityService(availabilityRepository).Upload(context.TODO(), data, service.HOSTED)

		assert.NoError(t, err)
		assert.Equal(t, response, url)
	})
}
