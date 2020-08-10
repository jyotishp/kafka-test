package http

import (
	"context"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/jyotishp/kafka-test/pkg/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

// StartHTTP registers the handlers and creates clients for our services.
func StartHTTP(port, grpcPort int) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.Dial(
		"localhost:"+strconv.Itoa(grpcPort),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(prometheus.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	rmux := runtime.NewServeMux()

	tweetsClient := pb.NewTweetsServiceClient(conn)
	err = pb.RegisterTweetsServiceHandlerClient(ctx, rmux, tweetsClient)
	if err != nil {
		log.Fatalf("failed to add gRPC proxy: %v", err)
	}
	log.Println("Registered with gRPC...")

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), mux)
	if err != nil {
		log.Fatal(err)
	}
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger-ui/app.swagger.json")
}
