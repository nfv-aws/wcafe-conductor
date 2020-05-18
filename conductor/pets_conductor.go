package conductor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/nfv-aws/wcafe-api-controller/config"
	"github.com/nfv-aws/wcafe-api-controller/db"
	"github.com/nfv-aws/wcafe-api-controller/entity"
	"log"
)

// User is alias of entity.Pets struct
type Pet entity.Pet

var (
	pets_svc       *sqs.SQS
	aws_region     string
	pets_queue_url string
)

func PetsInit() *sqs.SQS {
	config.Configure()
	aws_region = config.C.SQS.Region
	pets_queue_url = config.C.SQS.Pets_Queue_Url
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(aws_region)}))
	pets_svc := sqs.New(sess)
	return pets_svc
}

func PetsReceiveMessage(pets_svc *sqs.SQS) error {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(pets_queue_url),
		// 一度に取得する最大メッセージ数。最大でも1まで。
		MaxNumberOfMessages: aws.Int64(1),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	resp, err := pets_svc.ReceiveMessage(params)

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
		if err := PetsDeleteMessage(pets_svc, m); err != nil {
			log.Println(err)
		}
	}

	return nil
}

// メッセージを削除する。
func PetsDeleteMessage(pets_svc *sqs.SQS, msg *sqs.Message) error {
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
	pets_svc := PetsInit()
	for {
		if err := PetsReceiveMessage(pets_svc); err != nil {
			log.Fatal(err)
		}
	}
}
