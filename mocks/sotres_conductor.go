package mocks

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	gin "github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	entity "github.com/nfv-aws/wcafe-api-controller/entity"
)

type MockSQSSvc struct {
	sqsiface.SQSAPI
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

func (m *MockSQSSvc) StoreReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &sqs.ReceiveMessageOutput{}, nil
}

func (m *MockSQSSvc) StoreDeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return &sqs.DeleteMessageOutput{}, nil
}

// Update mocks base method.
func (m *MockStoreService) Update(id string, c *gin.Context) (entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, c)
	ret0, _ := ret[0].(entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
