package mocks

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	//	gin "github.com/gin-gonic/gin"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	entity "github.com/nfv-aws/wcafe-api-controller/entity"
)

type MockSQSSvc struct {
	sqsiface.SQSAPI
	Resp sqs.ReceiveMessageOutput
}

// MockStoreService is a mock of StoreService interface.
type MockStoreService struct {
	ctrl     *gomock.Controller
	recorder *MockStoreServiceMockRecorder
}

// MockStoreServiceMockRecorder is the mock recorder for MockStoreService.
type MockStoreServiceMockRecorder struct {
	mock *MockStoreService
}

// NewMockStoreService creates a new mock instance.
func NewMockStoreService(ctrl *gomock.Controller) *MockStoreService {
	mock := &MockStoreService{ctrl: ctrl}
	mock.recorder = &MockStoreServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreService) EXPECT() *MockStoreServiceMockRecorder {
	return m.recorder
}
func (m *MockSQSSvc) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &m.Resp, nil
}

func (m *MockSQSSvc) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return &sqs.DeleteMessageOutput{}, nil
}

func UpdateMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	um, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, nil, err
	}
	return um, mock, nil
}

// Update mocks base method.
func (m *MockStoreService) ChangeStrongPoint(id string) (entity.Store, error) {

	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id)
	ret0, _ := ret[0].(entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
