package conductor

import (
	"context"
	"fmt"
	"log"
	"net"

	"encoding/json"
	"google.golang.org/grpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	"github.com/nfv-aws/wcafe-api-controller/entity"
	pb "github.com/nfv-aws/wcafe-conductor/protoc"
)

// User is alias of entity.Stores struct
type Supply entity.Supply

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedSuppliesServer
}

func (s *server) SupplyList(ctx context.Context, in *pb.SupplyRequest) (*pb.SupplyResponse, error) {
	log.Printf("Received: %v", in.GetTable())
	list := SupplyList(in.GetTable())
	res, err := json.Marshal(list)
	if err != nil {
		log.Panic(err)
	}

	return &pb.SupplyResponse{Message: string(res)}, nil
}

func SuppliesGetMessage() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSuppliesServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

var (
	dynamodb *dynamo.DB
)

func SupplyList(target_table string) []Supply {
	db := dynamo.New(session.New(), &aws.Config{
		//	Credentials: c,
		Region: aws.String("ap-northeast-1"),
	})
	table := db.Table(target_table)
	var supplies []Supply
	err := table.Scan().All(&supplies)
	if err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
	log.Println(supplies)
	return supplies
}