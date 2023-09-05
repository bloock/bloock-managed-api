// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/repository/certification_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	domain "bloock-managed-api/internal/domain"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCertificationRepository is a mock of CertificationRepository interface.
type MockCertificationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCertificationRepositoryMockRecorder
}

// MockCertificationRepositoryMockRecorder is the mock recorder for MockCertificationRepository.
type MockCertificationRepositoryMockRecorder struct {
	mock *MockCertificationRepository
}

// NewMockCertificationRepository creates a new mock instance.
func NewMockCertificationRepository(ctrl *gomock.Controller) *MockCertificationRepository {
	mock := &MockCertificationRepository{ctrl: ctrl}
	mock.recorder = &MockCertificationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCertificationRepository) EXPECT() *MockCertificationRepositoryMockRecorder {
	return m.recorder
}

// ExistCertificationByHash mocks base method.
func (m *MockCertificationRepository) ExistCertificationByHash(ctx context.Context, hash string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistCertificationByHash", ctx, hash)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExistCertificationByHash indicates an expected call of ExistCertificationByHash.
func (mr *MockCertificationRepositoryMockRecorder) ExistCertificationByHash(ctx, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistCertificationByHash", reflect.TypeOf((*MockCertificationRepository)(nil).ExistCertificationByHash), ctx, hash)
}

// GetCertificationsByAnchorID mocks base method.
func (m *MockCertificationRepository) GetCertificationsByAnchorID(ctx context.Context, anchor int) ([]domain.Certification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificationsByAnchorID", ctx, anchor)
	ret0, _ := ret[0].([]domain.Certification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCertificationsByAnchorID indicates an expected call of GetCertificationsByAnchorID.
func (mr *MockCertificationRepositoryMockRecorder) GetCertificationsByAnchorID(ctx, anchor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificationsByAnchorID", reflect.TypeOf((*MockCertificationRepository)(nil).GetCertificationsByAnchorID), ctx, anchor)
}

// SaveCertification mocks base method.
func (m *MockCertificationRepository) SaveCertification(ctx context.Context, certification domain.Certification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCertification", ctx, certification)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCertification indicates an expected call of SaveCertification.
func (mr *MockCertificationRepositoryMockRecorder) SaveCertification(ctx, certification interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCertification", reflect.TypeOf((*MockCertificationRepository)(nil).SaveCertification), ctx, certification)
}

// UpdateCertificationDataID mocks base method.
func (m *MockCertificationRepository) UpdateCertificationDataID(ctx context.Context, certification domain.Certification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCertificationDataID", ctx, certification)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCertificationDataID indicates an expected call of UpdateCertificationDataID.
func (mr *MockCertificationRepositoryMockRecorder) UpdateCertificationDataID(ctx, certification interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCertificationDataID", reflect.TypeOf((*MockCertificationRepository)(nil).UpdateCertificationDataID), ctx, certification)
}
