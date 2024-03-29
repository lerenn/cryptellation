// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package exchanges is a generated GoMock package.
package exchanges

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

// Infos mocks base method.
func (m *MockPort) Infos(ctx context.Context, name string) (exchange.Exchange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Infos", ctx, name)
	ret0, _ := ret[0].(exchange.Exchange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Infos indicates an expected call of Infos.
func (mr *MockPortMockRecorder) Infos(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infos", reflect.TypeOf((*MockPort)(nil).Infos), ctx, name)
}
