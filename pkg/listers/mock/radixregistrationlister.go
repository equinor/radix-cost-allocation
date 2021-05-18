// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/listers/radixregistrationlister.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	v1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockRadixRegistrationLister is a mock of RadixRegistrationLister interface.
type MockRadixRegistrationLister struct {
	ctrl     *gomock.Controller
	recorder *MockRadixRegistrationListerMockRecorder
}

// MockRadixRegistrationListerMockRecorder is the mock recorder for MockRadixRegistrationLister.
type MockRadixRegistrationListerMockRecorder struct {
	mock *MockRadixRegistrationLister
}

// NewMockRadixRegistrationLister creates a new mock instance.
func NewMockRadixRegistrationLister(ctrl *gomock.Controller) *MockRadixRegistrationLister {
	mock := &MockRadixRegistrationLister{ctrl: ctrl}
	mock.recorder = &MockRadixRegistrationListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRadixRegistrationLister) EXPECT() *MockRadixRegistrationListerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockRadixRegistrationLister) List() []*v1.RadixRegistration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*v1.RadixRegistration)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockRadixRegistrationListerMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRadixRegistrationLister)(nil).List))
}
