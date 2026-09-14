package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	GRPCServerAddress  = "localhost:50053"
	healthCheckTimeOut = 3 * time.Second
)

func main() {
	addr := flag.String("addr", GRPCServerAddress, "gRPC server address")

	flag.Parse()

	ok := healthCheck(*addr)
	if !ok {
		os.Exit(1)
	}
}

func healthCheck(addr string) bool {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Print(fmt.Errorf("creating client: %w", err))

		return false
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print(fmt.Errorf("closing connection: %w", err))
		}
	}()

	client := healthpb.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeOut)
	defer cancel()

	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		log.Print(fmt.Errorf("checking health: %w", err))

		return false
	}

	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		log.Print("unhealthy upstream")

		return false
	}

	return true
}
