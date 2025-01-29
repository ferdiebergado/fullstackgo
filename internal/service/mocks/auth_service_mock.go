// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ferdiebergado/fullstackgo/internal/service (interfaces: AuthService)
//
// Generated by this command:
//
//	mockgen -destination=mocks/auth_service_mock.go -package=mocks . AuthService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/ferdiebergado/fullstackgo/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
	isgomock struct{}
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// SignInUser mocks base method.
func (m *MockAuthService) SignInUser(ctx context.Context, params model.UserSignInParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInUser", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignInUser indicates an expected call of SignInUser.
func (mr *MockAuthServiceMockRecorder) SignInUser(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInUser", reflect.TypeOf((*MockAuthService)(nil).SignInUser), ctx, params)
}

// SignUpUser mocks base method.
func (m *MockAuthService) SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUpUser", ctx, params)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUpUser indicates an expected call of SignUpUser.
func (mr *MockAuthServiceMockRecorder) SignUpUser(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUpUser", reflect.TypeOf((*MockAuthService)(nil).SignUpUser), ctx, params)
}
