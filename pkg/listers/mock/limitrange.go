// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/listers/limitrange.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
)

// MockLimitRangeLister is a mock of LimitRangeLister interface.
type MockLimitRangeLister struct {
	ctrl     *gomock.Controller
	recorder *MockLimitRangeListerMockRecorder
}

// MockLimitRangeListerMockRecorder is the mock recorder for MockLimitRangeLister.
type MockLimitRangeListerMockRecorder struct {
	mock *MockLimitRangeLister
}

// NewMockLimitRangeLister creates a new mock instance.
func NewMockLimitRangeLister(ctrl *gomock.Controller) *MockLimitRangeLister {
	mock := &MockLimitRangeLister{ctrl: ctrl}
	mock.recorder = &MockLimitRangeListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLimitRangeLister) EXPECT() *MockLimitRangeListerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockLimitRangeLister) List() []*v1.LimitRange {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*v1.LimitRange)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockLimitRangeListerMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockLimitRangeLister)(nil).List))
}
