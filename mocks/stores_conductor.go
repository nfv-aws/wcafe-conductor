package mocks

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	//	gin "github.com/gin-gonic/gin"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
)

type MockSQSSvc struct {
	sqsiface.SQSAPI
	Resp sqs.ReceiveMessageOutput
}

// MockStoreConductor is a mock of StoreConductor interface.
type MockStoreConductor struct {
	ctrl     *gomock.Controller
	recorder *MockStoreConductorMockRecorder
}

// MockStoreConductorMockRecorder is the mock recorder for MockStoreConductor.
type MockStoreConductorMockRecorder struct {
	mock *MockStoreConductor
}

// NewMockStoreConductor creates a new mock instance.
func NewMockStoreConductor(ctrl *gomock.Controller) *MockStoreConductor {
	mock := &MockStoreConductor{ctrl: ctrl}
	mock.recorder = &MockStoreConductorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreConductor) EXPECT() *MockStoreConductorMockRecorder {
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
	if err != nil {
		return nil, nil, err
	}
	um, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, nil, err
	}
	return um, mock, nil
}
