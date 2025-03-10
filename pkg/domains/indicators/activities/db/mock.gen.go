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

// ReadSMAActivity mocks base method.
func (m *MockDB) ReadSMAActivity(ctx context.Context, params ReadSMAActivityParams) (ReadSMAActivityResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSMAActivity", ctx, params)
	ret0, _ := ret[0].(ReadSMAActivityResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSMAActivity indicates an expected call of ReadSMAActivity.
func (mr *MockDBMockRecorder) ReadSMAActivity(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSMAActivity", reflect.TypeOf((*MockDB)(nil).ReadSMAActivity), ctx, params)
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

// UpsertSMAActivity mocks base method.
func (m *MockDB) UpsertSMAActivity(ctx context.Context, params UpsertSMAActivityParams) (UpsertSMAActivityResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertSMAActivity", ctx, params)
	ret0, _ := ret[0].(UpsertSMAActivityResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertSMAActivity indicates an expected call of UpsertSMAActivity.
func (mr *MockDBMockRecorder) UpsertSMAActivity(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertSMAActivity", reflect.TypeOf((*MockDB)(nil).UpsertSMAActivity), ctx, params)
}
