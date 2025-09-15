# copacetic-mcp

A Model Context Protocol (MCP) server for automated container patching using [Copacetic](https://github.com/project-copacetic/copacetic).

## MCP Tools

This server provides the following Model Context Protocol (MCP) tools:

- **`workflow-guide`**: Get patching strategy guide on which Copacetic tools to use for different container patching scenarios
- **`scan-container`**: Scan container images for vulnerabilities using Trivy - creates vulnerability reports required for report-based patching
- **`patch-vulnerabilities`**: Patch container image vulnerabilities using a pre-generated vulnerability report from 'scan-container' tool (RECOMMENDED approach for vulnerability-based patching)
- **`patch-platforms`**: Patch specific container image platforms with Copa - patches only the specified platforms WITHOUT vulnerability scanning
- **`patch-comprehensive`**: Comprehensively patch all container image platforms with Copa - patches all available platforms WITHOUT vulnerability scanning

## Installation

### VSCode Setup

To use copacetic-mcp with VSCode and MCP-compatible tools, add the following configuration to your VSCode `settings.json`:

```json
{
  "mcp.servers": {
    "copacetic-mcp": {
      "command": "/path/to/copacetic-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

Replace `/path/to/copacetic-mcp-server` with the actual path to your copacetic-mcp server binary. You can download the latest release from the [Releases]() page.

### Docker option (run server from a container)

```jsonc
"copacetic-mcp-docker": {
  "command": "docker",
  "args": [
    "run",
    "--rm",
    "-i",
    "--mount",
    "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock",
    "--mount",
    "type=bind,source=${env:HOME}/.docker/config.json,target=/root/.docker/config.json",
    "ghcr.io/duffney/copacetic-mcp:latest"
  ],
  "env": {
    "DOCKER_HOST": "unix:///var/run/docker.sock"
  }
}
```

Notes:

- Mounting the Docker socket gives the container access to the host Docker daemon; this is required for Copacetic image operations but has security implications—only run trusted images.
- Mounting `${HOME}/.docker/config.json` allows the container to use your registry credentials for pulling/pushing images.

#### Alternative: Using with Claude Desktop

You can also configure copacetic-mcp for use with Claude Desktop by adding it to your MCP configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "copacetic-mcp": {
      "command": "/path/to/copacetic-mcp-server",
      "args": []
    }
  }
}
```

## Development

### Prerequisites

- Go 1.20 or later
- [Copacetic](https://github.com/project-copacetic/copacetic) CLI installed
- [Trivy](https://github.com/aquasecurity/trivy) installed
- Docker with virtualization emulation support (required for container operations)

#### Copacetic CLI Requirements

Copacetic requires the following to be installed and available:

- **Docker**: Container runtime for image operations
- **Buildkit**: Advanced build features (included with recent Docker versions)
- **Container Registry Access**: For pulling and pushing patched images

### Project Structure

```
copacetic-mcp/
├── main.go                     # Main MCP server entry point
├── cmd/client/main.go         # Test client
├── internal/
│   ├── mcp/                   # MCP server handlers and setup
│   ├── copa/                  # Copacetic command execution
│   ├── trivy/                 # Trivy vulnerability scanning
│   ├── types/                 # Shared type definitions
│   └── util/                  # Utility functions (multiplatform, etc.)
├── .goreleaser.yml            # GoReleaser configuration
├── .github/workflows/         # GitHub Actions workflows
│   ├── build.yml             # Build and test on every push/PR
│   └── release.yml           # Automated releases on tags
└── Makefile                   # Development tasks
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
