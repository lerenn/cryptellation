// Code generated by MockGen. DO NOT EDIT.
// Source: db.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	reflect "reflect"

	worker "go.temporal.io/sdk/worker"
	gomock "go.uber.org/mock/gomock"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// CreateCandlesticksActivity mocks base method.
func (m *MockDB) CreateCandlesticksActivity(ctx context.Context, params CreateCandlesticksActivityParams) (CreateCandlesticksActivityResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCandlesticksActivity", ctx, params)
	ret0, _ := ret[0].(CreateCandlesticksActivityResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCandlesticksActivity indicates an expected call of CreateCandlesticksActivity.
func (mr *MockDBMockRecorder) CreateCandlesticksActivity(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCandlesticksActivity", reflect.TypeOf((*MockDB)(nil).CreateCandlesticksActivity), ctx, params)
}

// DeleteCandlesticksActivity mocks base method.
func (m *MockDB) DeleteCandlesticksActivity(ctx context.Context, params DeleteCandlesticksActivityParams) (DeleteCandlesticksActivityResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCandlesticksActivity", ctx, params)
	ret0, _ := ret[0].(DeleteCandlesticksActivityResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCandlesticksActivity indicates an expected call of DeleteCandlesticksActivity.
func (mr *MockDBMockRecorder) DeleteCandlesticksActivity(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCandlesticksActivity", reflect.TypeOf((*MockDB)(nil).DeleteCandlesticksActivity), ctx, params)
}

// ReadCandlesticksActivity mocks base method.
func (m *MockDB) ReadCandlesticksActivity(ctx context.Context, params ReadCandlesticksActivityParams) (ReadCandlesticksActivityResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCandlesticksActivity", ctx, params)
	ret0, _ := ret[0].(ReadCandlesticksActivityResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCandlesticksActivity indicates an expected call of ReadCandlesticksActivity.
func (mr *MockDBMockRecorder) ReadCandlesticksActivity(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCandlesticksActivity", reflect.TypeOf((*MockDB)(nil).ReadCandlesticksActivity), ctx, params)
}

// Register mocks base method.
func (m *MockDB) Register(w worker.Worker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", w)
}

// Register indicates an expected call of Register.
func (mr *MockDBMockRecorder) Register(w interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockDB)(nil).Register), w)
}

// UpdateCandlesticksActivity mocks base method.
func (m *MockDB) UpdateCandlesticksActivity(ctx context.Context, params UpdateCandlesticksActivityParams) (UpdateCandlesticksActivityResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCandlesticksActivity", ctx, params)
	ret0, _ := ret[0].(UpdateCandlesticksActivityResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCandlesticksActivity indicates an expected call of UpdateCandlesticksActivity.
func (mr *MockDBMockRecorder) UpdateCandlesticksActivity(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCandlesticksActivity", reflect.TypeOf((*MockDB)(nil).UpdateCandlesticksActivity), ctx, params)
}