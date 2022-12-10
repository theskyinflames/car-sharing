package main

import (
	"context"
)

const srvPort = ":80"

func main() {
	Run(context.Background(), srvPort)
}
