// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	reflect "reflect"
	time "time"

	candlestick "github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
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

// CreateCandlesticks mocks base method.
func (m *MockPort) CreateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCandlesticks", ctx, cs)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCandlesticks indicates an expected call of CreateCandlesticks.
func (mr *MockPortMockRecorder) CreateCandlesticks(ctx, cs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCandlesticks", reflect.TypeOf((*MockPort)(nil).CreateCandlesticks), ctx, cs)
}

// DeleteCandlesticks mocks base method.
func (m *MockPort) DeleteCandlesticks(ctx context.Context, cs *candlestick.List) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCandlesticks", ctx, cs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCandlesticks indicates an expected call of DeleteCandlesticks.
func (mr *MockPortMockRecorder) DeleteCandlesticks(ctx, cs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCandlesticks", reflect.TypeOf((*MockPort)(nil).DeleteCandlesticks), ctx, cs)
}

// ReadCandlesticks mocks base method.
func (m *MockPort) ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCandlesticks", ctx, cs, start, end, limit)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadCandlesticks indicates an expected call of ReadCandlesticks.
func (mr *MockPortMockRecorder) ReadCandlesticks(ctx, cs, start, end, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCandlesticks", reflect.TypeOf((*MockPort)(nil).ReadCandlesticks), ctx, cs, start, end, limit)
}

// UpdateCandlesticks mocks base method.
func (m *MockPort) UpdateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCandlesticks", ctx, cs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCandlesticks indicates an expected call of UpdateCandlesticks.
func (mr *MockPortMockRecorder) UpdateCandlesticks(ctx, cs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCandlesticks", reflect.TypeOf((*MockPort)(nil).UpdateCandlesticks), ctx, cs)
}
