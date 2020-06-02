package mocks

import (
	//	"github.com/golang/mock/gomock"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	gin "github.com/gin-gonic/gin"
	entity "github.com/nfv-aws/wcafe-api-controller/entity"
)

type MockSQSSvc struct {
	sqsiface.SQSAPI
}

func (m *MockSQSSvc) StoreReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &sqs.ReceiveMessageOutput{
		Messages: 1,
	}, nil
}

func (m *MockSQSSvc) StoreDeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return &sqs.ReceiveMessageOutput{
		Messages: 1,
	}, nil
}

// Update mocks base method.
func (m *MockStoreService) Update(id string, c *gin.Context) (entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, c)
	ret0, _ := ret[0].(entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
