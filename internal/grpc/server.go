package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/user/vigilante/internal/storage"
	"google.golang.org/grpc"
)

type Server struct{ DB *storage.DB }

func Start(ctx context.Context, port string, db *storage.DB) error {
	_ = db
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port)); if err != nil { return err }
	s := grpc.NewServer()
	go func(){ <-ctx.Done(); s.GracefulStop() }()
	return s.Serve(lis)
}
