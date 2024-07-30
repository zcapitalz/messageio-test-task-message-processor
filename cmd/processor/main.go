package main

import (
	"message-processor/internal/config"
	"message-processor/internal/processor"
)

func main() {
	cfg := config.MustNewProcessorConfig()
	processor.Process(&cfg)
}
