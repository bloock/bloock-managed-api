// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/repository/encryption_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	key "github.com/bloock/bloock-sdk-go/v2/entity/key"
	record "github.com/bloock/bloock-sdk-go/v2/entity/record"
	gomock "github.com/golang/mock/gomock"
)

// MockEncryptionRepository is a mock of EncryptionRepository interface.
type MockEncryptionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEncryptionRepositoryMockRecorder
}

// MockEncryptionRepositoryMockRecorder is the mock recorder for MockEncryptionRepository.
type MockEncryptionRepositoryMockRecorder struct {
	mock *MockEncryptionRepository
}

// NewMockEncryptionRepository creates a new mock instance.
func NewMockEncryptionRepository(ctrl *gomock.Controller) *MockEncryptionRepository {
	mock := &MockEncryptionRepository{ctrl: ctrl}
	mock.recorder = &MockEncryptionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEncryptionRepository) EXPECT() *MockEncryptionRepositoryMockRecorder {
	return m.recorder
}

// EncryptAESWithLocalKey mocks base method.
func (m *MockEncryptionRepository) EncryptAESWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, key string) (*record.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptAESWithLocalKey", ctx, data, kty, key)
	ret0, _ := ret[0].(*record.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptAESWithLocalKey indicates an expected call of EncryptAESWithLocalKey.
func (mr *MockEncryptionRepositoryMockRecorder) EncryptAESWithLocalKey(ctx, data, kty, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptAESWithLocalKey", reflect.TypeOf((*MockEncryptionRepository)(nil).EncryptAESWithLocalKey), ctx, data, kty, key)
}

// EncryptAESWithManagedKey mocks base method.
func (m *MockEncryptionRepository) EncryptAESWithManagedKey(ctx context.Context, data []byte, kid string) (*record.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptAESWithManagedKey", ctx, data, kid)
	ret0, _ := ret[0].(*record.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptAESWithManagedKey indicates an expected call of EncryptAESWithManagedKey.
func (mr *MockEncryptionRepositoryMockRecorder) EncryptAESWithManagedKey(ctx, data, kid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptAESWithManagedKey", reflect.TypeOf((*MockEncryptionRepository)(nil).EncryptAESWithManagedKey), ctx, data, kid)
}

// EncryptRSAWithLocalKey mocks base method.
func (m *MockEncryptionRepository) EncryptRSAWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, publicKey string, privateKey *string) (*record.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptRSAWithLocalKey", ctx, data, kty, publicKey, privateKey)
	ret0, _ := ret[0].(*record.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptRSAWithLocalKey indicates an expected call of EncryptRSAWithLocalKey.
func (mr *MockEncryptionRepositoryMockRecorder) EncryptRSAWithLocalKey(ctx, data, kty, publicKey, privateKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptRSAWithLocalKey", reflect.TypeOf((*MockEncryptionRepository)(nil).EncryptRSAWithLocalKey), ctx, data, kty, publicKey, privateKey)
}

// EncryptRSAWithManagedKey mocks base method.
func (m *MockEncryptionRepository) EncryptRSAWithManagedKey(ctx context.Context, data []byte, kid string) (*record.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptRSAWithManagedKey", ctx, data, kid)
	ret0, _ := ret[0].(*record.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptRSAWithManagedKey indicates an expected call of EncryptRSAWithManagedKey.
func (mr *MockEncryptionRepositoryMockRecorder) EncryptRSAWithManagedKey(ctx, data, kid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptRSAWithManagedKey", reflect.TypeOf((*MockEncryptionRepository)(nil).EncryptRSAWithManagedKey), ctx, data, kid)
}
