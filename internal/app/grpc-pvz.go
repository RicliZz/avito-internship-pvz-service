package app

import (
	"context"
	grpcServices "github.com/RicliZz/avito-internship-pvz-service/internal/services/grpc"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/proto"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func RungRPC() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		logger.Logger.Infow("Unable to connect to database:",
			"error", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", os.Getenv("PVZ_GRPC_ADDR")+os.Getenv("PVZ_GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pvz_v1.RegisterPVZServiceServer(s, &grpcServices.GRPCserver{Db: conn})
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
