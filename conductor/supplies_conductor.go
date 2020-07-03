package conductor

import (
	"context"
	"fmt"
	"net"

	"encoding/json"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	"github.com/nfv-aws/wcafe-api-controller/config"
	"github.com/nfv-aws/wcafe-api-controller/entity"
	pb "github.com/nfv-aws/wcafe-conductor/protoc"
)

// User is alias of entity.Stores struct
type Supply entity.Supply

type server struct {
	pb.UnimplementedSuppliesServer
}

func (s *server) SupplyList(ctx context.Context, in *pb.SupplyListRequest) (*pb.SupplyResponse, error) {
	log.Debug("SupplyList Receive gRPC Message: " + in.GetTable())
	list := supplylist(in.GetTable())
	res, err := json.Marshal(list)
	if err != nil {
		log.Panic(err)
	}
	return &pb.SupplyResponse{Message: string(res)}, nil
}

func (s *server) SupplyCreate(ctx context.Context, in *pb.SupplyCreateRequest) (*pb.SupplyResponse, error) {
	log.Debug("SupplyList Receive gRPC Message1: " + in.GetTable())
	log.Debug("SupplyList Receive gRPC Message2: " + in.GetBody())
	supply := supplycreate(in.GetTable(), in.GetBody())
	res, err := json.Marshal(supply)
	if err != nil {
		log.Panic(err)
	}
	return &pb.SupplyResponse{Message: string(res)}, nil
}

func SuppliesGetMessage() {
	config.Configure()
	log.Debug("Setting gRPC listten port")
	var port = ":" + config.C.Conductor.Port
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

func Dynamo_Init() *dynamo.DB {
	config.Configure()
	log.Debug("Init DynamoDB")
	aws_region = config.C.DynamoDB.Region
	dynamodb := dynamo.New(session.New(), &aws.Config{
		Region: aws.String(aws_region),
	})
	return dynamodb
}

func supplylist(target_table string) []Supply {
	dynamodb := Dynamo_Init()
	log.Debug("GetSupplyList by DynamoDB")
	table := dynamodb.Table(target_table)
	var supplies []Supply
	err := table.Scan().All(&supplies)
	if err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
	log.Println(supplies)
	return supplies
}

func supplycreate(target_table string, body string) Supply {
	dynamodb := Dynamo_Init()
	log.Debug("CreateSupply by DynamoDB")
	table := dynamodb.Table(target_table)
	var supply Supply
	err := json.Unmarshal([]byte(body), &supply)
	if err != nil {
		panic(err.Error())
	}
	log.Debug("Put Data")
	err = table.Put(supply).Run()
	if err != nil {
		panic(err.Error())
	}
	return supply
}
