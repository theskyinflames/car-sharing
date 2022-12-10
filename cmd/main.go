package main

import (
	"context"

	"theskyinflames/car-sharing/internal/infra/server"
)

const srvPort = ":80"

func main() {
	server.Run(context.Background(), srvPort)
}
