package hellod

import (
	"context"

	"go-monorepo/appmodule/helloer"
	"go-monorepo/database"
	"go-monorepo/internal/rpc/hello"
)

// Server is grpc server wrapper.
type Server struct {
	hello.UnimplementedGreeterServer

	service *helloer.Service
}

// NewServer creates server.
func NewServer() *Server {
	db := database.GetDB(database.Default)

	return &Server{
		service: helloer.NewService(db),
	}
}

// SayHello calls service.
func (s *Server) SayHello(
	ctx context.Context,
	req *hello.HelloRequest,
) (*hello.HelloReply, error) {
	resp, err := s.service.SayHello(req.Name)
	if err != nil {
		return nil, err
	}

	return &hello.HelloReply{
		Message: resp,
	}, nil
}
