package main

import (
	"message-processor/internal/config"
	"message-processor/internal/server"
)

func main() {
	cfg := config.MustNewServerConfig()
	server.Serve(cfg)
}
