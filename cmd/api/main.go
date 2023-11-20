package main

import (
	"context"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"

	"go-monorepo/app/api"
	"go-monorepo/health"
)

// Main starts process in cli.
func Main(ctx context.Context, c *cli.Context) {
	go health.StartServer()

	server := api.Server{}
	server.Start(ctx, c.String("listen-addr"))
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
