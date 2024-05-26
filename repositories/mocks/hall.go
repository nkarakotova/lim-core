// Code generated by MockGen. DO NOT EDIT.
// Source: hall.go

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	context "context"
	models "github.com/nkarakotova/lim-core/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHallRepository is a mock of HallRepository interface.
type MockHallRepository struct {
	ctrl     *gomock.Controller
	recorder *MockHallRepositoryMockRecorder
}

// MockHallRepositoryMockRecorder is the mock recorder for MockHallRepository.
type MockHallRepositoryMockRecorder struct {
	mock *MockHallRepository
}

// NewMockHallRepository creates a new mock instance.
func NewMockHallRepository(ctrl *gomock.Controller) *MockHallRepository {
	mock := &MockHallRepository{ctrl: ctrl}
	mock.recorder = &MockHallRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHallRepository) EXPECT() *MockHallRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockHallRepository) Create(ctx context.Context, hall *models.Hall) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, hall)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockHallRepositoryMockRecorder) Create(ctx, hall interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockHallRepository)(nil).Create), ctx, hall)
}

// GetAll mocks base method.
func (m *MockHallRepository) GetAll(ctx context.Context) (map[uint64]models.Hall, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].(map[uint64]models.Hall)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockHallRepositoryMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockHallRepository)(nil).GetAll), ctx)
}

// GetByID mocks base method.
func (m *MockHallRepository) GetByID(ctx context.Context, id uint64) (*models.Hall, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Hall)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockHallRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockHallRepository)(nil).GetByID), ctx, id)
}

// GetByNumber mocks base method.
func (m *MockHallRepository) GetByNumber(ctx context.Context, number uint64) (*models.Hall, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNumber", ctx, number)
	ret0, _ := ret[0].(*models.Hall)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNumber indicates an expected call of GetByNumber.
func (mr *MockHallRepositoryMockRecorder) GetByNumber(ctx, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNumber", reflect.TypeOf((*MockHallRepository)(nil).GetByNumber), ctx, number)
}