package conductor

import (
	//	"regexp"
	"testing"
	"time"

	//	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	//	"github.com/stretchr/testify/assert"

	"github.com/nfv-aws/wcafe-api-controller/entity"
	"github.com/nfv-aws/wcafe-conductor/mocks"
)

var (
	id_a = "37577ea0-6ebb-4a76-891b-f241e7e5dc7b"
	id_b = "78089ea0-58s3-6adk-8943-dejop43891f9"

	s = entity.Store{
		Id:          "sa5bafac-b35c-4852-82ca-b272cd79f2f3",
		Name:        "store_controller_test",
		Tag:         "Board game",
		Address:     "Shinagawa",
		StrongPoint: "helpful",
		CreatedAt:   ct,
		UpdatedAt:   ut,
	}
	ct, ut = time.Now(), time.Now()
)

func TestStoresReceiveMessageOK(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	err := StoresReceiveMessage(svc)
	if err != nil {
		t.Error(err)
	}

}

// ** ToDo チケット158 **
//func TestStoresReceiveMessageOK2(t *testing.T) {
//	svc := &mocks.MockSQSSvc{}
//	svc.Resp = sqs.ReceiveMessageOutput{
//		Messages: []*sqs.Message{
//			{Body: aws.String(id_a)},
//			{Body: aws.String(id_b)},
//		}}
//	err := StoresReceiveMessage(svc)
//	if err != nil {
//		t.Error(err)
//	}
//
//}

func TestStoresDeleteMessageOK(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	msg := &sqs.Message{
		ReceiptHandle: aws.String("delete")}
	err := StoresDeleteMessage(svc, msg)
	if err != nil {
		t.Error(err)
	}
}

// ** ToDo チケット158**
// DB ExpectQueryの返り値が戻ってこない
//func TestChangeStrongPointOK(t *testing.T) {
//	db, mock, err := mocks.UpdateMock()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer db.Close()
//	db.LogMode(true)

//	m := &sqs.Message{
//		Body: aws.String(id_a)}

//	e := entity.Store{
//		Id:          *m.Body,
//		Name:        "store_controller_test",
//		Tag:         "Board game",
//		Address:     "Shinagawa",
//		StrongPoint: "sqs_test",
//		CreatedAt:   s.CreatedAt,
//		UpdatedAt:   s.UpdatedAt,
//	}

//	mock.ExpectQuery(regexp.QuoteMeta(
//		"UPDATE stores SET strong_point='sqs_test', updated_at='2020-06-03 13:20:57' WHERE id='37577ea0-6ebb-4a76-891b-f241e7e5dc7b';")).
//		WithArgs("id").WillReturnRows(
//		sqlmock.NewRows([]string{"id", "name", "tag", "address", "strong_point", "created_at", "updated_at"}).
//			AddRow(*m.Body, s.Name, s.Tag, s.Address, "strong_point", s.CreatedAt, s.UpdatedAt))

//	resp, err := ChangeStrongPoint(*m.Body, db)
//	if err != nil {
//		t.Error(err)
//	}
//	assert.Equal(t, e, resp)

//}
