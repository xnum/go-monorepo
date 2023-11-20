package main

import (
	"context"
	"log"
	"net"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"go-monorepo/app/hellod"
	"go-monorepo/health"
	"go-monorepo/internal/rpc/hello"
)

// Main starts process in cli.
func Main(ctx context.Context, c *cli.Context) {
	go health.StartServer()

	server := hellod.NewServer()

	lis, err := net.Listen("tcp", c.String("listen-addr"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	app := boot.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen-addr",
				Value: ":8787",
			},
		},
		Main: Main,
	}

	app.Run()
}
