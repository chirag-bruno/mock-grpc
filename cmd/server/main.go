package main

import (
	"log"
	"os"

	"github.com/chirag-bruno/mock-grpc/internal/server"
	"github.com/chirag-bruno/mock-grpc/internal/transport"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "todo-server",
		Usage: "A gRPC server for todo management",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Value:   "http",
				Usage:   "Server mode: http, unix, or pipe",
			},
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Value:   "localhost:50051",
				Usage:   "Server address (for http: host:port, for unix: socket path)",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	mode, err := transport.ParseMode(c.String("mode"))
	if err != nil {
		return err
	}

	return server.Run(server.Config{
		Mode:    mode,
		Address: c.String("address"),
	})
}
