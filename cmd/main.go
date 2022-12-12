package main

import (
	"context"

	"theskyinflames/car-sharing/cmd/service"
)

const srvPort = ":80"

func main() {
	service.Run(context.Background(), srvPort)
}
