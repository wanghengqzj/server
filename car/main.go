package main

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"server/car/mq/amqpclt"
	"server/car/sim"
	"server/car/sim/pos"
	"server/car/trip"
	"server/car/ws"
	rentalpb "server/rental/api/gen/v1"
	coolenvpb "server/shared/coolenv"

	//"net/http"
	carpb "server/car/api/gen/v1"
	"server/car/car"
	"server/car/dao"
	//"server/car/mq/amqpclt"
	//"server/car/sim"
	//"server/car/sim/pos"
	//"server/car/trip"
	//"server/car/ws"
	//rentalpb "server/rental/api/gen/v1"
	//coolenvpb "server/shared/coolenv"
	"server/shared/server"

	//"github.com/gorilla/websocket"
	"github.com/namsral/flag"
	//"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8084", "address to listen")
var wsAddr = flag.String("ws_addr", ":9090", "websocket address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://root:root@localhost:27017", "mongo uri")
var amqpURL = flag.String("amqp_url", "amqp://guest:guest@localhost:5672/", "amqp url")
var carAddr = flag.String("car_addr", "localhost:8084", "address for car service")
var tripAddr = flag.String("trip_addr", "localhost:8082", "address for trip service")
var aiAddr = flag.String("ai_addr", "localhost:18001", "address for ai service")

func main() {
	//flag.Parse()

	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	//mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(*mongoURI))
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	db := mongoClient.Database("test")

	amqpConn, err := amqp.Dial(*amqpURL)
	if err != nil {
		logger.Fatal("cannot dial amqp", zap.Error(err))
	}

	exchange := "coolcar"
	pub, err := amqpclt.NewPublisher(amqpConn, exchange)
	if err != nil {
		logger.Fatal("cannot create publisher", zap.Error(err))
	}
	//
	//// Run car simulations.
	//carConn, err := grpc.Dial(*carAddr, grpc.WithInsecure())
	carConn, err := grpc.Dial("localhost:8084", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect car service", zap.Error(err))
	}
	aiConn, err := grpc.Dial("localhost:18001", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect ai service", zap.Error(err))
	}
	sub, err := amqpclt.NewSubscriber(amqpConn, exchange, logger)
	if err != nil {
		logger.Fatal("cannot create subscriber", zap.Error(err))
	}
	posSub, err := amqpclt.NewSubscriber(amqpConn, "pos_sim", logger)
	if err != nil {
		logger.Fatal("cannot create pos subscriber", zap.Error(err))
	}
	simController := &sim.Controller{
		CarService:    carpb.NewCarServiceClient(carConn),
		AIService:     coolenvpb.NewAIServiceClient(aiConn),
		Logger:        logger,
		CarSubscriber: sub,
		PosSubscriber: &pos.Subscriber{
			Sub:    posSub,
			Logger: logger,
		},
	}
	go simController.RunSimulations(context.Background())

	// Start websocket handler.
	u := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc("/ws", ws.Handler(u, sub, logger))
	go func() {
		//addr := *wsAddr
		addr := ":9090"
		logger.Info("HTTP server started.", zap.String("addr", addr))
		logger.Sugar().Fatal(
			http.ListenAndServe(addr, nil))
	}()
	// Start trip updater.
	tripConn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect trip service", zap.Error(err))
	}
	go trip.RunUpdater(sub, rentalpb.NewTripServiceClient(tripConn), logger)

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name: "car",
		//Addr:   *addr,
		Addr:   ":8084",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			carpb.RegisterCarServiceServer(s, &car.Service{
				Logger:    logger,
				Mongo:     dao.NewMongo(db),
				Publisher: pub,
			})
		},
	}))
}
