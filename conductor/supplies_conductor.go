package conductor

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

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
	log.Debug("SupplyCreate Receive gRPC Message1: " + in.GetTable())
	log.Debug("SupplyCreate Receive gRPC Message2: " + in.GetBody())
	supply := supplycreate(in.GetTable(), in.GetBody())
	res, err := json.Marshal(supply)
	if err != nil {
		log.Panic(err)
	}
	return &pb.SupplyResponse{Message: string(res)}, nil
}

func (s *server) SupplyUpdate(ctx context.Context, in *pb.SupplyUpdateRequest) (*pb.SupplyResponse, error) {
	log.Debug("SupplyUpdate Receive gRPC Message1: " + in.GetTable())
	log.Debug("SupplyUpdate Receive gRPC Message2: " + in.GetId())
	log.Debug("SupplyUpdate Receive gRPC Message3: " + in.GetBody())

	supply, err := supplyupdate(in.GetTable(), in.GetId(), in.GetBody())
	if err != nil {
		log.Error(err)
		return &pb.SupplyResponse{}, err
	}
	res, err := json.Marshal(supply)
	if err != nil {
		log.Panic(err)
		return &pb.SupplyResponse{Message: string(res)}, err
	}
	return &pb.SupplyResponse{Message: string(res)}, nil
}

func (s *server) SupplyDelete(ctx context.Context, in *pb.SupplyDeleteRequest) (*pb.SupplyResponse, error) {
	log.Debug("SupplyDelete Receive gRPC Message1: " + in.GetTable())
	log.Debug("SupplyDelete Receive gRPC Message2: " + in.GetId())
	supply, err := supplydelete(in.GetTable(), in.GetId())
	if err != nil {
		return nil, err
	}
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

func supplyupdate(target_table string, id string, body string) (Supply, error) {
	dynamodb := Dynamo_Init()
	log.Debug("UpdateSupply by DynamoDB")
	table := dynamodb.Table(target_table)
	var supply Supply

	if err := table.Get("id", id).One(&supply); err != nil {
		return supply, err
	}

	err := json.Unmarshal([]byte(body), &supply)
	if err != nil {
		panic(err.Error())
		return supply, err
	}

	log.Debug("Update Data")
	err = table.Update("id", id).Set("name", supply.Name).Set("price", supply.Price).Set("type", supply.Type).Value(&supply)
	if err != nil {
		panic(err.Error())
		return supply, err
	}
	return supply, nil
}

func supplydelete(target_table string, id string) (Supply, error) {
	dynamodb := Dynamo_Init()
	log.Debug("DeleteSupply by DynamoDB")
	table := dynamodb.Table(target_table)
	var supply Supply

	if err := table.Get("id", id).One(&supply); err != nil {
		return supply, err
	}

	log.Debug("Delete Data")
	if err := table.Delete("id", id).Run(); err != nil {
		panic(err.Error())
	}
	return supply, nil
}
