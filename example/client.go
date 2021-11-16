package main

import (
	"context"
	"fmt"
	"log"
	"time"

	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"

	opaV1 "github.com/gmhafiz/opa_service/api/v1"
)

type grpcClient struct {
	grpcClient opaV1.ServiceClient
}

// Run the gRPC server first, then run this client example with
// go run cmd/opa/main.go
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	opts := []grpcRetry.CallOption{
		grpcRetry.WithBackoff(
			grpcRetry.BackoffLinearWithJitter(100*time.Millisecond, 10),
		),
	}
	conn, err := grpc.DialContext(ctx, ":9091",
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpcRetry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpcRetry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	client := grpcClient{}
	client.grpcClient = opaV1.NewServiceClient(conn)

	liveness, err := client.grpcClient.Liveness(ctx, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(liveness)

	req := &opaV1.CheckRequest{
		User:     1,
		Resource: "referees",
		Action:   "GET",
	}

	checkResponse, err := client.grpcClient.IsAllowed(ctx, req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(checkResponse)
}
