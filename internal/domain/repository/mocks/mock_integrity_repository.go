// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/repository/integrity_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	domain "bloock-managed-api/internal/domain"
	context "context"
	reflect "reflect"

	integrity "github.com/bloock/bloock-sdk-go/v2/entity/integrity"
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
func (m *MockIntegrityRepository) Certify(ctx context.Context, bytes [][]byte) ([]domain.Certification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Certify", ctx, bytes)
	ret0, _ := ret[0].([]domain.Certification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Certify indicates an expected call of Certify.
func (mr *MockIntegrityRepositoryMockRecorder) Certify(ctx, bytes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Certify", reflect.TypeOf((*MockIntegrityRepository)(nil).Certify), ctx, bytes)
}

// GetAnchorByID mocks base method.
func (m *MockIntegrityRepository) GetAnchorByID(ctx context.Context, anchorID int) (integrity.Anchor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnchorByID", ctx, anchorID)
	ret0, _ := ret[0].(integrity.Anchor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnchorByID indicates an expected call of GetAnchorByID.
func (mr *MockIntegrityRepositoryMockRecorder) GetAnchorByID(ctx, anchorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnchorByID", reflect.TypeOf((*MockIntegrityRepository)(nil).GetAnchorByID), ctx, anchorID)
}
