package main

import (
	"net"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	pb "github.com/teralion/live-connections/disk/proto"
)

const (
	port = 50051
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

func (s *areaManagerServer) Create(ctx context.Context, request *pb.CreateAreaRequest) (*pb.CreateAreaResponse, error) {
	log.Println("Received create request")
	return &pb.CreateAreaResponse{}, nil
}

func (s *areaManagerServer) ListUsers(ctx context.Context, request *pb.ListAreaUsersRequest) (*pb.ListAreaUsersResponse, error) {
	log.Println("Received list users request")
	return &pb.ListAreaUsersResponse{}, nil
}

func (s *areaManagerServer) Destroy(ctx context.Context, request *pb.DestroyAreaRequest) (*pb.DestroyAreaResponse, error) {
	log.Println("Received destroy request")
	return &pb.DestroyAreaResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAreaManagerServer(grpcServer, &areaManagerServer{})
	pb.RegisterUserManagerServer(grpcServer, &userManagerServer{})
	pb.RegisterShapeManagerServer(grpcServer, &shapeManagerServer{})
	log.Printf("server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}