package main

import (
	"context"
	"log"

	"github.com/duffney/copacetic-mcp/internal/mcp"
)

// Build information set by GoReleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := mcp.Run(context.Background(), version); err != nil {
		log.Fatal(err)
	}
}
