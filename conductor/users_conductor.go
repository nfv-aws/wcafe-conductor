package conductor

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/jinzhu/gorm"

	"github.com/nfv-aws/wcafe-api-controller/config"
	"github.com/nfv-aws/wcafe-api-controller/db"
	"github.com/nfv-aws/wcafe-api-controller/entity"
)

// User is alias of entity.Users struct
type User entity.User

var (
	users_svc       *sqs.SQS
	users_queue_url string
)

func UsersInit() *sqs.SQS {
	log.Debug("Init Users")
	config.Configure()
	aws_region = config.C.SQS.Region
	users_queue_url = config.C.SQS.Users_Queue_Url
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(aws_region)}))
	users_svc := sqs.New(sess)
	return users_svc
}

func UsersReceiveMessage(users_svc sqsiface.SQSAPI) (*sqs.ReceiveMessageOutput, error) {
	log.Debug("UsersReceiveMessage")
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(users_queue_url),
		// 一度に取得する最大メッセージ数。最大でも1まで。
		MaxNumberOfMessages: aws.Int64(1),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	resp, err := users_svc.ReceiveMessage(params)

	if err != nil {
		return resp, err
	}

	log.WithFields(log.Fields{
		"count": len(resp.Messages),
	}).Info("messages ")

	// 取得したキューの数が0の場合emptyと表示
	if len(resp.Messages) == 0 {
		log.Info("empty queue.")
	}

	return resp, nil
}

func UsersChangeDB(users_svc sqsiface.SQSAPI, resp *sqs.ReceiveMessageOutput) error {
	log.Debug("UsersChangeDB")
	db := db.GetDB()
	// メッセージの数だけループを回し、userのStatusを変更する
	for _, m := range resp.Messages {
		log.Debug(*m.Body)
		if err := UsersChangeStatus(*m.Body, db); err != nil {
			log.Fatal(err)
			return err
		}
		// 処理が終わったキューを削除
		if err := UsersDeleteMessage(users_svc, m); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// メッセージを削除する。
func UsersDeleteMessage(users_svc sqsiface.SQSAPI, msg *sqs.Message) error {
	log.Debug("UsersDeleteMessageDe")
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(users_queue_url),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := users_svc.DeleteMessage(params)

	if err != nil {
		return err
	}
	return nil
}

// DBのStatusをCREATEに変更する
func UsersChangeStatus(id string, db *gorm.DB) error {
	log.Debug("UsersChangeStatus")
	var u entity.User

	// usersのstatusを変更
	u.Status = "CREATED"

	if err := db.Table("users").Where("id = ?", id).Updates(&u).Error; err != nil {
		return err
	}
	log.Println("CHANGE Users Status")
	return nil
}

// キューを刈り取り、usersのPOST時の処理をおこなう
func UsersGetMessage() {
	log.Debug("UsersGetMessage")
	users_svc := UsersInit()
	for {
		resp, err := UsersReceiveMessage(users_svc)
		if err != nil {
			log.Fatal(err)
		}
		if err := UsersChangeDB(users_svc, resp); err != nil {
			log.Fatal(err)
		}
	}
}
