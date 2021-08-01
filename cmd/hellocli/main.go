package main

import (
	"context"
	"log"

	"github.com/urfave/cli/v2"

	"go-monorepo/app"
	"go-monorepo/internal/rpc/hello"
)

var (
	clientConfig = hello.NewClientConfig("")
	name         string
)

// Main starts process in cli.
func Main(ctx context.Context, c *cli.Context) {
	clientPool := hello.NewClientPoolFromConfig(clientConfig)
	msg, err := clientPool.SayHello(ctx, name)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(msg)
}

func main() {
	app := app.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Required:    true,
				Destination: &name,
			},
		},
		Main: Main,
	}

	app.Run()
}
