// Code generated by MockGen. DO NOT EDIT.
// Source: repository.mailpicker.mock.go

// Package mock is a generated GoMock package.
package MailPicker

import (
	model "2019_2_Next_Level/internal/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// UserExists mocks base method
func (m *MockRepository) UserExists(login string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExists", login)
	ret0, _ := ret[0].(bool)
	return ret0
}

// UserExists indicates an expected call of UserExists
func (mr *MockRepositoryMockRecorder) UserExists(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExists", reflect.TypeOf((*MockRepository)(nil).UserExists), login)
}

// AddEmail mocks base method
func (m *MockRepository) AddEmail(arg0 *model.Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEmail", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEmail indicates an expected call of AddEmail
func (mr *MockRepositoryMockRecorder) AddEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEmail", reflect.TypeOf((*MockRepository)(nil).AddEmail), arg0)
}
