package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	pb "github.com/codex-team/hawk.cloud-manager/api/protobuf"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedConfigServiceServer
}

func (s *server) Get(ctx context.Context, in *empty.Empty) (*pb.Config, error) {
	return &pb.Config{}, nil
}

func Run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
