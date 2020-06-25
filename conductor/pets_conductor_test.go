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

func TestPetsReceiveMessageOK(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	resp, err := PetsReceiveMessage(svc)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)
}

func TestPetsReceiveMessageOK2(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	svc.Resp = sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			{Body: aws.String(id_a)},
			{Body: aws.String(id_b)},
		}}
	resp, err := PetsReceiveMessage(svc)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)
}

func TestPetsDeleteMessageOK(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	msg := &sqs.Message{
		ReceiptHandle: aws.String("delete")}
	err := PetsDeleteMessage(svc, msg)
	if err != nil {
		t.Error(err)
	}
}

func TestPetsChangeStatusOK(t *testing.T) {
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
		"UPDATE `pets` SET `status` = ? WHERE (id = ?)")).
		WithArgs("CREATED", *m.Body).WillReturnResult(
		sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := PetsChangeStatus(*m.Body, db); err != nil {
		t.Error(err)
	}
}
