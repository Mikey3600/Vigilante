package grpc

import (
	"fmt"
	"net"

	"github.com/user/vigilante/internal/storage"
	"google.golang.org/grpc"
)

// Server implements the generated MetricIngestionServer interface.
type Server struct {
	DB *storage.DB
}

// Start opens up a gRPC bounding listener.
func Start(port string, db *storage.DB) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	// When code is compiled, RegisterMetricIngestionServer goes here.
	return s.Serve(lis)
}
