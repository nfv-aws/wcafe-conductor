package conductor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/nfv-aws/wcafe-api-controller/config"
	"github.com/nfv-aws/wcafe-api-controller/db"
	"github.com/nfv-aws/wcafe-api-controller/entity"
	"log"
)

// User is alias of entity.Pets struct
type Pet entity.Pet

type Queue struct {
	Client sqsiface.SQSAPI
	URL    string
}

var (
	aws_region string
)

func PetsInit() Queue {
	config.Configure()
	aws_region = config.C.SQS.Region
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(aws_region)}))
	q := Queue{
		Client: sqs.New(sess),
		URL:    config.C.SQS.Pets_Queue_Url,
	}
	return q
}

func (q *Queue) PetsReceiveMessage() error {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(q.URL),
		// 一度に取得する最大メッセージ数。最大でも1まで。
		MaxNumberOfMessages: aws.Int64(1),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	resp, err := q.Client.ReceiveMessage(params)

	if err != nil {
		return err
	}

	log.Printf("messages count: %d\n", len(resp.Messages))

	// 取得したキューの数が0の場合emptyと表示
	if len(resp.Messages) == 0 {
		log.Println("empty queue.")
		return nil
	}

	// メッセージの数だけループを回し、petのSTATUSの値を変更する
	for _, m := range resp.Messages {
		log.Println(*m.Body)
		ChangeStatus(*m.Body)
		// 処理が終わったキューを削除
		if err := q.PetsDeleteMessage(m); err != nil {
			log.Println(err)
		}
	}

	return nil
}

// メッセージを削除する。
func (q *Queue) PetsDeleteMessage(msg *sqs.Message) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.URL),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := q.Client.DeleteMessage(params)

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

// キューを刈り取り、petsのPOST時の処理をおこなう
func PetsGetMessage() {
	q := PetsInit()

	for {
		if err := q.PetsReceiveMessage(); err != nil {
			log.Fatal(err)
		}
	}
}
