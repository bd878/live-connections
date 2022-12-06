package main

import (
	"net"
	"fmt"
	"log"

	"google.golang.org/grpc"
	pb "github.com/teralion/live-connections/proto"
)

const (
	port = 5051
)

type areaManagerServer struct {
	pb.UnimplementedAreaManagerServer
}

type userManagerServer struct {
	pb.UnimplementedUserManagerServer
}

type shapeManagerServer struct {
	pb.UnimplementedShapeManagerServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAreaManagerService(grpcServer, &areaManagerServer{})
	pb.RegisterUserManagerService(grpcServer, &userManagerServer{})
	pb.RegisterShapeManagerService(grpcServer, &shapeManagerServer{})
	log.Printf("server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}