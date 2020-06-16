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

// User is alias of entity.Stores struct
type Store entity.Store

var (
	stores_svc       *sqs.SQS
	stores_queue_url string
)

func StoresInit() *sqs.SQS {
	log.Debug("Init Stores")
	config.Configure()
	aws_region = config.C.SQS.Region
	stores_queue_url = config.C.SQS.Stores_Queue_Url
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(aws_region)}))
	stores_svc := sqs.New(sess)
	return stores_svc
}

func StoresReceiveMessage(stores_svc sqsiface.SQSAPI) error {
	log.Debug("StoresReceiveMessage")
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(stores_queue_url),
		// 一度に取得する最大メッセージ数。最大でも1まで。
		MaxNumberOfMessages: aws.Int64(1),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	resp, err := stores_svc.ReceiveMessage(params)

	if err != nil {
		return resp, err
	}

	log.Printf("messages count: %d\n", len(resp.Messages))

	// 取得したキューの数が0の場合emptyと表示
	if len(resp.Messages) == 0 {
		log.Println("empty queue.")
	}

	return resp, nil
}

func StoresChangeDB(stores_svc sqsiface.SQSAPI, resp *sqs.ReceiveMessageOutput) error {
	log.Debug("StoresChangeDB")
	db := db.GetDB()
	// メッセージの数だけループを回し、storeのStrongPointを変更する
	for _, m := range resp.Messages {
		log.Println(*m.Body)
		ChangeStrongPoint(*m.Body, db)
		// 処理が終わったキューを削除
		if err := StoresDeleteMessage(stores_svc, m); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// メッセージを削除する。
func StoresDeleteMessage(stores_svc sqsiface.SQSAPI, msg *sqs.Message) error {
	log.Debug("StoresDeleteMessage")
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(stores_queue_url),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := stores_svc.DeleteMessage(params)

	if err != nil {
		return err
	}
	return nil
}

// DBのStrongPointを"sqs_test"に変更する
func ChangeStrongPoint(id string, db *gorm.DB) (entity.Store, error) {
	log.Debug("ChangeStrongPoint")
	var u entity.Store

	// storesのStrongPointを変更
	u.StrongPoint = "sqs_test"

	if err := db.Table("stores").Where("id = ?", id).Updates(&u).Error; err != nil {
		return u, err
	}
	log.Println("CHANGE StrongPoint")
	return u, nil
}

// キューを刈り取り、storesのPOST時の処理をおこなう
func StoresGetMessage() {
	log.Debug("StoresGetMessage")
	stores_svc := StoresInit()
	for {
		resp, err := StoresReceiveMessage(stores_svc)
		if err != nil {
			log.Fatal(err)
		}
		if err := StoresChangeDB(stores_svc, resp); err != nil {
			log.Fatal(err)
		}
	}

}
