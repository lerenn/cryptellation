// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	reflect "reflect"

	exchange "github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
	gomock "go.uber.org/mock/gomock"
)

// MockPort is a mock of Port interface.
type MockPort struct {
	ctrl     *gomock.Controller
	recorder *MockPortMockRecorder
}

// MockPortMockRecorder is the mock recorder for MockPort.
type MockPortMockRecorder struct {
	mock *MockPort
}

// NewMockPort creates a new mock instance.
func NewMockPort(ctrl *gomock.Controller) *MockPort {
	mock := &MockPort{ctrl: ctrl}
	mock.recorder = &MockPortMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPort) EXPECT() *MockPortMockRecorder {
	return m.recorder
}

// CreateExchanges mocks base method.
func (m *MockPort) CreateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range exchanges {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateExchanges", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateExchanges indicates an expected call of CreateExchanges.
func (mr *MockPortMockRecorder) CreateExchanges(ctx interface{}, exchanges ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, exchanges...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExchanges", reflect.TypeOf((*MockPort)(nil).CreateExchanges), varargs...)
}

// DeleteExchanges mocks base method.
func (m *MockPort) DeleteExchanges(ctx context.Context, names ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range names {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteExchanges", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExchanges indicates an expected call of DeleteExchanges.
func (mr *MockPortMockRecorder) DeleteExchanges(ctx interface{}, names ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, names...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExchanges", reflect.TypeOf((*MockPort)(nil).DeleteExchanges), varargs...)
}

// ReadExchanges mocks base method.
func (m *MockPort) ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range names {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadExchanges", varargs...)
	ret0, _ := ret[0].([]exchange.Exchange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadExchanges indicates an expected call of ReadExchanges.
func (mr *MockPortMockRecorder) ReadExchanges(ctx interface{}, names ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, names...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadExchanges", reflect.TypeOf((*MockPort)(nil).ReadExchanges), varargs...)
}

// UpdateExchanges mocks base method.
func (m *MockPort) UpdateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range exchanges {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateExchanges", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExchanges indicates an expected call of UpdateExchanges.
func (mr *MockPortMockRecorder) UpdateExchanges(ctx interface{}, exchanges ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, exchanges...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExchanges", reflect.TypeOf((*MockPort)(nil).UpdateExchanges), varargs...)
}
