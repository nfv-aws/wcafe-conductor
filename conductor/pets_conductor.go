package conductor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/nfv-aws/wcafe-conductor/config"
	"github.com/nfv-aws/wcafe-conductor/db"
	"github.com/nfv-aws/wcafe-conductor/entity"
	"log"
)

// User is alias of entity.Pets struct
type Pet entity.Pet

var (
	svc        *sqs.SQS
	aws_region string
	queue_url  string
)

func Init() *sqs.SQS {
	config.Configure()
	aws_region = config.C.SQS.Region
	queue_url = config.C.SQS.Queue_Url
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(aws_region)}))
	svc := sqs.New(sess)
	return svc
}

func ReceiveMessage(svc *sqs.SQS) error {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queue_url),
		// 一度に取得する最大メッセージ数。最大でも1まで。
		MaxNumberOfMessages: aws.Int64(1),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		return err
	}

	log.Printf("messages count: %d\n", len(resp.Messages))

	// 取得したキューの数が0の場合emptyと表示
	if len(resp.Messages) == 0 {
		log.Println("empty queue.")
		return nil
	}

	// メッセージの数だけループを回し、STATUSの値を変更
	for _, m := range resp.Messages {
		log.Println(*m.Body)
		ChangeStatus(*m.Body)
		// 処理が終わったキューを削除
		if err := DeleteMessage(svc, m); err != nil {
			log.Println(err)
		}
	}

	return nil
}

// メッセージを削除する。
func DeleteMessage(svc *sqs.SQS, msg *sqs.Message) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue_url),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := svc.DeleteMessage(params)

	if err != nil {
		return err
	}
	return nil
}

// DBのStatusをCREATEに変更する
func ChangeStatus(id string) (Pet, error) {
	db := db.GetDB()
	var u Pet

	// petsのstatusを変更
	u.Status = "CREATED"

	if err := db.Table("pets").Where("id = ?", id).Updates(&u).Error; err != nil {
		return u, err
	}
	log.Println("CHANGE STATUS")

	return u, nil
}

// キューを刈り取り、POSTの処理のSTATUSの値を"CREATED"に書き換える
func GetMessage() {
	svc := Init()
	for {
		if err := ReceiveMessage(svc); err != nil {
			log.Fatal(err)
		}
	}
}
