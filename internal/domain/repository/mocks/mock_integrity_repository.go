// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/repository/integrity_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	domain "github.com/bloock/bloock-managed-api/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockIntegrityRepository is a mock of IntegrityRepository interface.
type MockIntegrityRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIntegrityRepositoryMockRecorder
}

// MockIntegrityRepositoryMockRecorder is the mock recorder for MockIntegrityRepository.
type MockIntegrityRepositoryMockRecorder struct {
	mock *MockIntegrityRepository
}

// NewMockIntegrityRepository creates a new mock instance.
func NewMockIntegrityRepository(ctrl *gomock.Controller) *MockIntegrityRepository {
	mock := &MockIntegrityRepository{ctrl: ctrl}
	mock.recorder = &MockIntegrityRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIntegrityRepository) EXPECT() *MockIntegrityRepositoryMockRecorder {
	return m.recorder
}

// Certify mocks base method.
func (m *MockIntegrityRepository) Certify(ctx context.Context, file []byte) (domain.Certification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Certify", ctx, file)
	ret0, _ := ret[0].(domain.Certification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Certify indicates an expected call of Certify.
func (mr *MockIntegrityRepositoryMockRecorder) Certify(ctx, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Certify", reflect.TypeOf((*MockIntegrityRepository)(nil).Certify), ctx, file)
}
