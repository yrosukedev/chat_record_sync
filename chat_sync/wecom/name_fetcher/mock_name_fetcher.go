// Code generated by MockGen. DO NOT EDIT.
// Source: chat_sync/wecom/transformer/name_fetcher.go

// Package name_fetcher is a generated GoMock package.
package name_fetcher

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockNameFetcher is a mock of NameFetcher interface.
type MockNameFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockNameFetcherMockRecorder
}

// MockNameFetcherMockRecorder is the mock recorder for MockNameFetcher.
type MockNameFetcherMockRecorder struct {
	mock *MockNameFetcher
}

// NewMockNameFetcher creates a new mock instance.
func NewMockNameFetcher(ctrl *gomock.Controller) *MockNameFetcher {
	mock := &MockNameFetcher{ctrl: ctrl}
	mock.recorder = &MockNameFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNameFetcher) EXPECT() *MockNameFetcherMockRecorder {
	return m.recorder
}

// FetchName mocks base method.
func (m *MockNameFetcher) FetchName(id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchName", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchName indicates an expected call of FetchName.
func (mr *MockNameFetcherMockRecorder) FetchName(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchName", reflect.TypeOf((*MockNameFetcher)(nil).FetchName), id)
}
