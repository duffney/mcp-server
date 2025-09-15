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

/*
	TODO: Patch all local images and overwrite current tag
	TODO: Scan tool, return vulns wit sev.
	TODO: Run mcp server from a container to avoid having to install/config tools
	TODO: Integrate with contagious and patch entire registry
*/

func main() {
	if err := mcp.Run(context.Background(), version); err != nil {
		log.Fatal(err)
	}
}
