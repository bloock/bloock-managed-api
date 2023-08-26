// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/application_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	service "bloock-managed-api/internal/service"
	request "bloock-managed-api/internal/service/authenticity/request"
	response "bloock-managed-api/internal/service/authenticity/response"
	request0 "bloock-managed-api/internal/service/integrity/request"
	response0 "bloock-managed-api/internal/service/integrity/response"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBaseProcessService is a mock of BaseProcessService interface.
type MockBaseProcessService struct {
	ctrl     *gomock.Controller
	recorder *MockBaseProcessServiceMockRecorder
}

// MockBaseProcessServiceMockRecorder is the mock recorder for MockBaseProcessService.
type MockBaseProcessServiceMockRecorder struct {
	mock *MockBaseProcessService
}

// NewMockBaseProcessService creates a new mock instance.
func NewMockBaseProcessService(ctrl *gomock.Controller) *MockBaseProcessService {
	mock := &MockBaseProcessService{ctrl: ctrl}
	mock.recorder = &MockBaseProcessServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBaseProcessService) EXPECT() *MockBaseProcessServiceMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockBaseProcessService) Process(ctx context.Context, req service.ProcessRequest) (*response.ProcessResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", ctx, req)
	ret0, _ := ret[0].(*response.ProcessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockBaseProcessServiceMockRecorder) Process(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockBaseProcessService)(nil).Process), ctx, req)
}

// MockAuthenticityService is a mock of AuthenticityService interface.
type MockAuthenticityService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticityServiceMockRecorder
}

// MockAuthenticityServiceMockRecorder is the mock recorder for MockAuthenticityService.
type MockAuthenticityServiceMockRecorder struct {
	mock *MockAuthenticityService
}

// NewMockAuthenticityService creates a new mock instance.
func NewMockAuthenticityService(ctrl *gomock.Controller) *MockAuthenticityService {
	mock := &MockAuthenticityService{ctrl: ctrl}
	mock.recorder = &MockAuthenticityServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticityService) EXPECT() *MockAuthenticityServiceMockRecorder {
	return m.recorder
}

// Sign mocks base method.
func (m *MockAuthenticityService) Sign(ctx context.Context, SignRequest request.SignRequest) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", ctx, SignRequest)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Sign indicates an expected call of Sign.
func (mr *MockAuthenticityServiceMockRecorder) Sign(ctx, SignRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockAuthenticityService)(nil).Sign), ctx, SignRequest)
}

// MockIntegrityService is a mock of IntegrityService interface.
type MockIntegrityService struct {
	ctrl     *gomock.Controller
	recorder *MockIntegrityServiceMockRecorder
}

// MockIntegrityServiceMockRecorder is the mock recorder for MockIntegrityService.
type MockIntegrityServiceMockRecorder struct {
	mock *MockIntegrityService
}

// NewMockIntegrityService creates a new mock instance.
func NewMockIntegrityService(ctrl *gomock.Controller) *MockIntegrityService {
	mock := &MockIntegrityService{ctrl: ctrl}
	mock.recorder = &MockIntegrityServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIntegrityService) EXPECT() *MockIntegrityServiceMockRecorder {
	return m.recorder
}

// Certify mocks base method.
func (m *MockIntegrityService) Certify(ctx context.Context, files []byte) ([]response0.CertificationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Certify", ctx, files)
	ret0, _ := ret[0].([]response0.CertificationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Certify indicates an expected call of Certify.
func (mr *MockIntegrityServiceMockRecorder) Certify(ctx, files interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Certify", reflect.TypeOf((*MockIntegrityService)(nil).Certify), ctx, files)
}

// MockAvailabilityService is a mock of AvailabilityService interface.
type MockAvailabilityService struct {
	ctrl     *gomock.Controller
	recorder *MockAvailabilityServiceMockRecorder
}

// MockAvailabilityServiceMockRecorder is the mock recorder for MockAvailabilityService.
type MockAvailabilityServiceMockRecorder struct {
	mock *MockAvailabilityService
}

// NewMockAvailabilityService creates a new mock instance.
func NewMockAvailabilityService(ctrl *gomock.Controller) *MockAvailabilityService {
	mock := &MockAvailabilityService{ctrl: ctrl}
	mock.recorder = &MockAvailabilityServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAvailabilityService) EXPECT() *MockAvailabilityServiceMockRecorder {
	return m.recorder
}

// UploadHosted mocks base method.
func (m *MockAvailabilityService) UploadHosted(ctx context.Context, data []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadHosted", ctx, data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadHosted indicates an expected call of UploadHosted.
func (mr *MockAvailabilityServiceMockRecorder) UploadHosted(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadHosted", reflect.TypeOf((*MockAvailabilityService)(nil).UploadHosted), ctx, data)
}

// UploadIpfs mocks base method.
func (m *MockAvailabilityService) UploadIpfs(ctx context.Context, data []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadIpfs", ctx, data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadIpfs indicates an expected call of UploadIpfs.
func (mr *MockAvailabilityServiceMockRecorder) UploadIpfs(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadIpfs", reflect.TypeOf((*MockAvailabilityService)(nil).UploadIpfs), ctx, data)
}

// MockCertificateUpdateAnchorService is a mock of CertificateUpdateAnchorService interface.
type MockCertificateUpdateAnchorService struct {
	ctrl     *gomock.Controller
	recorder *MockCertificateUpdateAnchorServiceMockRecorder
}

// MockCertificateUpdateAnchorServiceMockRecorder is the mock recorder for MockCertificateUpdateAnchorService.
type MockCertificateUpdateAnchorServiceMockRecorder struct {
	mock *MockCertificateUpdateAnchorService
}

// NewMockCertificateUpdateAnchorService creates a new mock instance.
func NewMockCertificateUpdateAnchorService(ctrl *gomock.Controller) *MockCertificateUpdateAnchorService {
	mock := &MockCertificateUpdateAnchorService{ctrl: ctrl}
	mock.recorder = &MockCertificateUpdateAnchorServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCertificateUpdateAnchorService) EXPECT() *MockCertificateUpdateAnchorServiceMockRecorder {
	return m.recorder
}

// UpdateAnchor mocks base method.
func (m *MockCertificateUpdateAnchorService) UpdateAnchor(ctx context.Context, updateRequest request0.UpdateCertificationAnchorRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAnchor", ctx, updateRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAnchor indicates an expected call of UpdateAnchor.
func (mr *MockCertificateUpdateAnchorServiceMockRecorder) UpdateAnchor(ctx, updateRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAnchor", reflect.TypeOf((*MockCertificateUpdateAnchorService)(nil).UpdateAnchor), ctx, updateRequest)
}
