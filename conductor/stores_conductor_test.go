package conductor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/nfv-aws/wcafe-conductor/mocks"
	//	"reflect"
	//	gin "github.com/gin-gonic/gin"
	"github.com/nfv-aws/wcafe-api-controller/entity"
	"github.com/stretchr/testify/assert"
	//	"net/http/httptest"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"

	"testing"
	"time"
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

func TestChangeStrongPointOK(t *testing.T) {
	db, mock, err := mocks.UpdateMock()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true)

	m := &sqs.Message{
		Body: aws.String(id_a)}

	r := entity.Store{Id: *m.Body}

	mock.ExpectQuery(regexp.QuoteMeta(
		"UPDATE stores SET strong_point='sqs_test', updated_at=time.Now() WHERE id=*m.Body AND active=true;")).
		WithArgs("id").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "tag", "address", "strong_point", "created_at", "updated_at"}).
			AddRow(r.Id))

	resp, err := ChangeStrongPoint(*m.Body)
	if err != nil {
		t.Error(err)
	}
	assert.ElementsMatch(t, r, resp)

}
