package mocks

import (
	"github.com/golang/mock/gomock"
	"marketplace-api/internal/models"
	"reflect"
)

//go:generate mockgen -source=../service.go -destination=mocks/listing_service_mock.go

type MockListingService struct {
	ctrl     *gomock.Controller
	recorder *MockListingServiceMockRecorder
}

type MockListingServiceMockRecorder struct {
	mock *MockListingService
}

func NewMockListingService(ctrl *gomock.Controller) *MockListingService {
	mock := &MockListingService{ctrl: ctrl}
	mock.recorder = &MockListingServiceMockRecorder{mock}
	return mock
}

func (m *MockListingService) EXPECT() *MockListingServiceMockRecorder {
	return m.recorder
}

func (m *MockListingService) CreateListing(userID int, req models.CreateListingRequest) (*models.Listing, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateListing", userID, req)
	ret0, _ := ret[0].(*models.Listing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockListingServiceMockRecorder) CreateListing(userID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateListing", reflect.TypeOf((*MockListingService)(nil).CreateListing), userID, req)
}

func (m *MockListingService) GetListings(filter models.ListingsFilter, currentUserID *int) (*models.PaginatedListings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListings", filter, currentUserID)
	ret0, _ := ret[0].(*models.PaginatedListings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockListingServiceMockRecorder) GetListings(filter, currentUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListings", reflect.TypeOf((*MockListingService)(nil).GetListings), filter, currentUserID)
}

func (m *MockListingService) GetListingByID(id int, currentUserID *int) (*models.Listing, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListingByID", id, currentUserID)
	ret0, _ := ret[0].(*models.Listing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockListingServiceMockRecorder) GetListingByID(id, currentUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListingByID", reflect.TypeOf((*MockListingService)(nil).GetListingByID), id, currentUserID)
}

func (m *MockListingService) UpdateListing(id int, userID int, req models.UpdateListingRequest) (*models.Listing, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateListing", id, userID, req)
	ret0, _ := ret[0].(*models.Listing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockListingServiceMockRecorder) UpdateListing(id, userID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateListing", reflect.TypeOf((*MockListingService)(nil).UpdateListing), id, userID, req)
}

func (m *MockListingService) DeleteListing(id int, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteListing", id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockListingServiceMockRecorder) DeleteListing(id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteListing", reflect.TypeOf((*MockListingService)(nil).DeleteListing), id, userID)
}

func (m *MockListingService) GetUserListings(userID int, filter models.ListingsFilter) (*models.PaginatedListings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserListings", userID, filter)
	ret0, _ := ret[0].(*models.PaginatedListings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockListingServiceMockRecorder) GetUserListings(userID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserListings", reflect.TypeOf((*MockListingService)(nil).GetUserListings), userID, filter)
}
