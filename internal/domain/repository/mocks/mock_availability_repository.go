// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/repository/availability_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	domain "github.com/bloock/bloock-managed-api/internal/domain"
	record "github.com/bloock/bloock-sdk-go/v2/entity/record"
	gomock "github.com/golang/mock/gomock"
)

// MockAvailabilityRepository is a mock of AvailabilityRepository interface.
type MockAvailabilityRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAvailabilityRepositoryMockRecorder
}

// MockAvailabilityRepositoryMockRecorder is the mock recorder for MockAvailabilityRepository.
type MockAvailabilityRepositoryMockRecorder struct {
	mock *MockAvailabilityRepository
}

// NewMockAvailabilityRepository creates a new mock instance.
func NewMockAvailabilityRepository(ctrl *gomock.Controller) *MockAvailabilityRepository {
	mock := &MockAvailabilityRepository{ctrl: ctrl}
	mock.recorder = &MockAvailabilityRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAvailabilityRepository) EXPECT() *MockAvailabilityRepositoryMockRecorder {
	return m.recorder
}

// FindFile mocks base method.
func (m *MockAvailabilityRepository) FindFile(ctx context.Context, dataID string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindFile", ctx, dataID)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFile indicates an expected call of FindFile.
func (mr *MockAvailabilityRepositoryMockRecorder) FindFile(ctx, dataID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFile", reflect.TypeOf((*MockAvailabilityRepository)(nil).FindFile), ctx, dataID)
}

// RetrieveLocal mocks base method.
func (m *MockAvailabilityRepository) RetrieveLocal(ctx context.Context, filePath string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveLocal", ctx, filePath)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveLocal indicates an expected call of RetrieveLocal.
func (mr *MockAvailabilityRepositoryMockRecorder) RetrieveLocal(ctx, filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveLocal", reflect.TypeOf((*MockAvailabilityRepository)(nil).RetrieveLocal), ctx, filePath)
}

// RetrieveTmp mocks base method.
func (m *MockAvailabilityRepository) RetrieveTmp(ctx context.Context, filename string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveTmp", ctx, filename)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveTmp indicates an expected call of RetrieveTmp.
func (mr *MockAvailabilityRepositoryMockRecorder) RetrieveTmp(ctx, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveTmp", reflect.TypeOf((*MockAvailabilityRepository)(nil).RetrieveTmp), ctx, filename)
}

// UploadHosted mocks base method.
func (m *MockAvailabilityRepository) UploadHosted(ctx context.Context, file *domain.File, record record.Record) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadHosted", ctx, file, record)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadHosted indicates an expected call of UploadHosted.
func (mr *MockAvailabilityRepositoryMockRecorder) UploadHosted(ctx, file, record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadHosted", reflect.TypeOf((*MockAvailabilityRepository)(nil).UploadHosted), ctx, file, record)
}

// UploadIpfs mocks base method.
func (m *MockAvailabilityRepository) UploadIpfs(ctx context.Context, file *domain.File, record record.Record) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadIpfs", ctx, file, record)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadIpfs indicates an expected call of UploadIpfs.
func (mr *MockAvailabilityRepositoryMockRecorder) UploadIpfs(ctx, file, record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadIpfs", reflect.TypeOf((*MockAvailabilityRepository)(nil).UploadIpfs), ctx, file, record)
}

// UploadLocal mocks base method.
func (m *MockAvailabilityRepository) UploadLocal(ctx context.Context, file *domain.File) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadLocal", ctx, file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadLocal indicates an expected call of UploadLocal.
func (mr *MockAvailabilityRepositoryMockRecorder) UploadLocal(ctx, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadLocal", reflect.TypeOf((*MockAvailabilityRepository)(nil).UploadLocal), ctx, file)
}

// UploadTmp mocks base method.
func (m *MockAvailabilityRepository) UploadTmp(ctx context.Context, file *domain.File, record record.Record) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadTmp", ctx, file, record)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadTmp indicates an expected call of UploadTmp.
func (mr *MockAvailabilityRepositoryMockRecorder) UploadTmp(ctx, file, record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadTmp", reflect.TypeOf((*MockAvailabilityRepository)(nil).UploadTmp), ctx, file, record)
}
