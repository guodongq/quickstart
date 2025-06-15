package commands

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"

	pb "github.com/guodongq/quickstart/pkg/api/genproto/helloworld/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               "go-grpc-quickstart",
		Short:             "A gRPC server hello world example",
		DisableAutoGenTag: true,
		Run: func(c *cobra.Command, _ []string) {
			flag.Parse()
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			reflection.Register(s)

			grpc_health_v1.RegisterHealthServer(s, health.NewServer())
			pb.RegisterGreeterServer(s, &server{})
			log.Printf("server listening at %v", lis.Addr())
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		},
	}

	return command
}

var port = flag.Int("port", 8080, "The server port")

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
