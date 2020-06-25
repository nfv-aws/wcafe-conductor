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

func TestUsersReceiveMessageOK(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	resp, err := UsersReceiveMessage(svc)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)
}

func TestUsersReceiveMessageOK2(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	svc.Resp = sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			{Body: aws.String(id_a)},
			{Body: aws.String(id_b)},
		}}
	resp, err := UsersReceiveMessage(svc)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)
}

func TestUsersDeleteMessageOK(t *testing.T) {
	svc := &mocks.MockSQSSvc{}
	msg := &sqs.Message{
		ReceiptHandle: aws.String("delete")}
	err := UsersDeleteMessage(svc, msg)
	if err != nil {
		t.Error(err)
	}
}

func TestUsersChangeStatusOK(t *testing.T) {
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
		"UPDATE `users` SET `address` = ? WHERE (id = ?)")).
		WithArgs("Kyoto", *m.Body).WillReturnResult(
		sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := UsersChangeStatus(*m.Body, db); err != nil {
		t.Error(err)
	}
}
