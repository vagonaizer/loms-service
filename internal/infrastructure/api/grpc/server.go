package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/vagonaizer/loms/api/protos/gen/loms"
	"github.com/vagonaizer/loms/internal/usecase/loms"
)

type Server struct {
	server *grpc.Server
	port   int
}

func NewServer(port int, service *loms.Service) *Server {
	log.Printf("Initializing gRPC server on port %d", port)
	server := grpc.NewServer()
	pb.RegisterLOMSServer(server, service)

	return &Server{
		server: server,
		port:   port,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting gRPC server on %s", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("Server is listening on %s", addr)
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) Stop() {
	log.Println("Stopping gRPC server...")
	s.server.GracefulStop()
	log.Println("gRPC server stopped")
}
