package http

import (
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/jyotishp/kafka-test/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

// Start gRPC server
func StartGRPC(port int, kafkaTopic string, brokers []string) {
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to bind address: %v", err)
	}

	grpcServer := grpc.NewServer()
	tweetsServer := NewTweetsServer(kafkaTopic, brokers)
	pb.RegisterTweetsServiceServer(grpcServer, tweetsServer)
	prometheus.Register(grpcServer)
	log.Println("gRPC server ready...")
	grpcServer.Serve(lis)
}
