package main

import (
	"context"

	"github.com/reiot777/protobuf-example/server"
)

func main() {
	var srv = server.Server{
		Host: "0.0.0.0",
		Port: 8888,
	}

	srv.Serve(context.Background())
}
