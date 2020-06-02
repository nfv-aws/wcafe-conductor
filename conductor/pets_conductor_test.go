package conductor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/assert"
	"testing"
	//	"github.com/nfv-aws/wcafe-conductor/mocks"
)

type MockedReceiveMsgs struct {
	sqsiface.SQSAPI
	Resp sqs.ReceiveMessageOutput
}

func (m MockedReceiveMsgs) ReceiveMessage(in *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	// Only need to return mocked response output
	return &m.Resp, nil
}

func TestQueueGetMessage(t *testing.T) {
	resp := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			Body: aws.String("74684838-a5d9-47d8-91a4-ff63ce802763")},
	}

	q := Queue{
		Client: mockedReceiveMsgs{Resp: c.Resp},
		URL:    "https://USERS-QUEUE",
	}
	err := q.PetsReceiveMessage()
	assert.Equal(t, "74684838-a5d9-47d8-91a4-ff63ce802763", []resp.Messages.Body)
}
