package main

import (
	"context"
	"log"

	"github.com/project-copacetic/mcp-server/internal/copamcp"
)

// Build information set by GoReleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := copamcp.Run(context.Background(), version); err != nil {
		log.Fatal(err)
	}
}
