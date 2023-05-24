package main

import (
	"context"
	"fmt"
	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	blobpb "server/blob/api/gen/v1"
	"server/blob/blob"
	"server/blob/cos"
	"server/blob/dao"
	"server/shared/server"
)

var addr = flag.String("addr", ":8083", "address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://root:root@localhost:27017", "mongo uri")
var cosAddr = flag.String("cos_addr", "<URL>", "cos address")
var cosSecID = flag.String("cos_sec_id", "<SEC_ID>", "cos secret id")
var cosSecKey = flag.String("cos_sec_key", "<SEC_KEY>", "cos secret key")

func main() {
	//flag.Parse()

	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	db := mongoClient.Database("test")

	//st, err := cos.NewService(*cosAddr, *cosSecID, *cosSecKey)
	st, err := cos.NewService(
		"https://coolcarlearn-1317438012.cos.ap-beijing.myqcloud.com",
		"AKIDHVDmDQPDtSyVRXju12tdhcMvJXYMyQoK",
		"XvAKSi087xfGkMsIZlprofO5e2JRYYN6")
	fmt.Println("cos service...")
	if err != nil {
		fmt.Println("cannot create cos service....")
		logger.Fatal("cannot create cos service", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name: "blob",
		//Addr:   *addr,
		Addr:   ":8083",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			blobpb.RegisterBlobServiceServer(s, &blob.Service{
				Storage: st,
				Mongo:   dao.NewMongo(db),
				Logger:  logger,
			})
		},
	}))
}
