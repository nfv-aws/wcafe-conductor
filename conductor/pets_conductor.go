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

// User is alias of entity.Pets struct
type Pet entity.Pet

var (
	aws_region     string
	pets_svc       *sqs.SQS
	pets_queue_url string
)

func PetsInit() *sqs.SQS {
	log.Debug("Init Pets")
	config.Configure()
	aws_region = config.C.SQS.Region
	pets_queue_url = config.C.SQS.Pets_Queue_Url
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(aws_region)}))
	pets_svc := sqs.New(sess)
	return pets_svc
}

func PetsReceiveMessage(pets_svc sqsiface.SQSAPI) (*sqs.ReceiveMessageOutput, error) {
	log.Debug("PetsReceiveMessage")
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(pets_queue_url),
		// 一度に取得する最大メッセージ数。最大でも1まで。
		MaxNumberOfMessages: aws.Int64(1),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	resp, err := pets_svc.ReceiveMessage(params)

	if err != nil {
		return resp, err
	}

	log.Info("messages count: " + string(len(resp.Messages)) + "\n")

	// 取得したキューの数が0の場合emptyと表示
	if len(resp.Messages) == 0 {
		log.Info("empty queue.")
	}

	return resp, nil
}

func PetsChangeDB(pets_svc sqsiface.SQSAPI, resp *sqs.ReceiveMessageOutput) error {
	log.Debug("PetsChangeDB")
	db := db.GetDB()
	// メッセージの数だけループを回し、petのStrongPointを変更する
	for _, m := range resp.Messages {
		log.Debug(*m.Body)
		if err := PetsChangeStatus(*m.Body, db); err != nil {
			log.Fatal(err)
			return err
		}
		// 処理が終わったキューを削除
		if err := PetsDeleteMessage(pets_svc, m); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// メッセージを削除する。
func PetsDeleteMessage(pets_svc sqsiface.SQSAPI, msg *sqs.Message) error {
	log.Debug("PetsDeleteMessage")
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(pets_queue_url),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := pets_svc.DeleteMessage(params)

	if err != nil {
		return err
	}
	return nil
}

// DBのStatusをCREATEに変更する
func PetsChangeStatus(id string, db *gorm.DB) error {
	log.Debug("PetsChangeStatus")
	var u entity.Pet

	// petsのstatusを変更
	u.Status = "CREATED"

	if err := db.Table("pets").Where("id = ?", id).Updates(&u).Error; err != nil {
		return err
	}
	log.Println("CHANGE Status")
	return nil
}

// キューを刈り取り、petsのPOST時の処理をおこなう
func PetsGetMessage() {
	log.Debug("PetsGetMessage")
	pets_svc := PetsInit()
	for {
		resp, err := PetsReceiveMessage(pets_svc)
		if err != nil {
			log.Fatal(err)
		}
		if err := PetsChangeDB(pets_svc, resp); err != nil {
			log.Fatal(err)
		}
	}
}
