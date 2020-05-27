package conductor

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

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
	config.Configure()
	aws_region = config.C.SQS.Region
	users_queue_url = config.C.SQS.Users_Queue_Url
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(aws_region)}))
	users_svc := sqs.New(sess)
	return users_svc
}

func UsersReceiveMessage(users_svc *sqs.SQS) error {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(users_queue_url),
		// 一度に取得する最大メッセージ数。最大でも1まで。
		MaxNumberOfMessages: aws.Int64(1),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	resp, err := users_svc.ReceiveMessage(params)

	if err != nil {
		return err
	}

	log.Printf("messages count: %d\n", len(resp.Messages))

	// 取得したキューの数が0の場合emptyと表示
	if len(resp.Messages) == 0 {
		log.Println("empty queue.")
		return nil
	}

	// メッセージの数だけループを回し、userのSTATUSの値を変更する
	for _, m := range resp.Messages {
		log.Println(*m.Body)
		ChangeUserStatus(*m.Body)
		// 処理が終わったキューを削除
		if err := UsersDeleteMessage(users_svc, m); err != nil {
			log.Println(err)
		}
	}
	return nil
}

// メッセージを削除する。
func UsersDeleteMessage(users_svc *sqs.SQS, msg *sqs.Message) error {
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
func ChangeUserStatus(id string) (User, error) {
	db := db.GetDB()
	var u User

	// usersのstatusを変更(今はお試しでAddressを変更)
	u.Address = "Kyoto"

	if err := db.Table("users").Where("id = ?", id).Updates(&u).Error; err != nil {
		return u, err
	}
	log.Println("CHANGE STATUS")
	return u, nil
}

// キューを刈り取り、usersのPOST時の処理をおこなう
func UsersGetMessage() {
	users_svc := UsersInit()
	for {
		if err := UsersReceiveMessage(users_svc); err != nil {
			log.Fatal(err)
		}
	}
}
