package mocks

import (
	"github.com/golang/mock/gomock"
	"marketplace-api/internal/models"
	"reflect"
)

//go:generate mockgen -source=../service.go -destination=mocks/auth_service_mock.go

type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

func (m *MockAuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", req)
	ret0, _ := ret[0].(*models.AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceMockRecorder) Register(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthService)(nil).Register), req)
}

func (m *MockAuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", req)
	ret0, _ := ret[0].(*models.AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceMockRecorder) Login(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), req)
}

func (m *MockAuthService) GetUserByID(id int) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockAuthService)(nil).GetUserByID), id)
}
