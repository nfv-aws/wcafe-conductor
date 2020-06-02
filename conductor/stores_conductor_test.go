package conductor

import (
	//	"github.com/aws/aws-sdk-go/aws"
	//	"github.com/aws/aws-sdk-go/service/sqs"
	//	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/nfv-aws/wcafe-conductor/mocks"
	"reflect"
	"testing"
)

func TestStoreReceiveMessageOK(t *testing.T) {
	svc := &MockSQSSvc
	r, err := mock.StoreReceiveMessage(svc)
	if err != nil {
		t.Error(err)

		expected := 1
	}
	if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %q to eq %q", expected, r)
	}
}
