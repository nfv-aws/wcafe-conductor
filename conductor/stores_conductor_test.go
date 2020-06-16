package conductor

import (
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/nfv-aws/wcafe-conductor/mocks"
)

var (
	id_a = "37577ea0-6ebb-4a76-891b-f241e7e5dc7b"
	id_b = "78089ea0-58s3-6adk-8943-dejop43891f9"
)

func TestStoresReceiveMessageOK(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	resp, err := StoresReceiveMessage(svc)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)
}

func TestStoresReceiveMessageOK2(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	svc.Resp = sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			{Body: aws.String(id_a)},
			{Body: aws.String(id_b)},
		}}
	resp, err := StoresReceiveMessage(svc)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)
}

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
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE `stores` SET `strong_point` = ? WHERE (id = ?)")).
		WithArgs("sqs_test", *m.Body).WillReturnResult(
		sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := ChangeStrongPoint(*m.Body, db); err != nil {
		t.Error(err)
	}
}
