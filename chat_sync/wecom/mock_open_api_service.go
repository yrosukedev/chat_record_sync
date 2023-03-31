// Code generated by MockGen. DO NOT EDIT.
// Source: chat_sync/wecom/open_api_service.go

// Package wecom is a generated GoMock package.
package wecom

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOpenAPIService is a mock of OpenAPIService interface.
type MockOpenAPIService struct {
	ctrl     *gomock.Controller
	recorder *MockOpenAPIServiceMockRecorder
}

// MockOpenAPIServiceMockRecorder is the mock recorder for MockOpenAPIService.
type MockOpenAPIServiceMockRecorder struct {
	mock *MockOpenAPIService
}

// NewMockOpenAPIService creates a new mock instance.
func NewMockOpenAPIService(ctrl *gomock.Controller) *MockOpenAPIService {
	mock := &MockOpenAPIService{ctrl: ctrl}
	mock.recorder = &MockOpenAPIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOpenAPIService) EXPECT() *MockOpenAPIServiceMockRecorder {
	return m.recorder
}

// GetExternalContactByID mocks base method.
func (m *MockOpenAPIService) GetExternalContactByID(externalId string) (*ExternalContact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExternalContactByID", externalId)
	ret0, _ := ret[0].(*ExternalContact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExternalContactByID indicates an expected call of GetExternalContactByID.
func (mr *MockOpenAPIServiceMockRecorder) GetExternalContactByID(externalId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExternalContactByID", reflect.TypeOf((*MockOpenAPIService)(nil).GetExternalContactByID), externalId)
}

// GetUserInfoByID mocks base method.
func (m *MockOpenAPIService) GetUserInfoByID(id string) (*UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfoByID", id)
	ret0, _ := ret[0].(*UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfoByID indicates an expected call of GetUserInfoByID.
func (mr *MockOpenAPIServiceMockRecorder) GetUserInfoByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfoByID", reflect.TypeOf((*MockOpenAPIService)(nil).GetUserInfoByID), id)
}