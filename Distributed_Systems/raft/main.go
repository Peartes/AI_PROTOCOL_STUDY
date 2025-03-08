package main

import (
	"log"
	"net"

	"github.com/peartes/distr_system/raft/rpc"
	types "github.com/peartes/distr_system/raft/types"
	"google.golang.org/grpc"
)

func main() {
	// Create a listener on a TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the RaftService server
	types.RegisterRaftServerServer(grpcServer, &rpc.RaftServer{})

	log.Println("gRPC server listening on port 50051")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
