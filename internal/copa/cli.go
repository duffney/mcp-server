package copa

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/project-copacetic/mcp-server/internal/docker"
	"github.com/project-copacetic/mcp-server/internal/types"
)

const (
	defaultVexFile = "vex.json"
)

// CopaSupportedPlatforms lists all platforms that Copa can patch
// Based on Copa documentation: https://project-copacetic.github.io/copacetic/website/multiplatform-patching
var CopaSupportedPlatforms = []string{
	"linux/amd64",
	"linux/arm64",
	"linux/arm/v7",
	"linux/arm/v6",
	"linux/386",
	"linux/ppc64le",
	"linux/s390x",
	"linux/riscv64",
}

type CLI struct {
	copaPath   string
	dryRun     bool
	image      string
	tag        string
	platforms  []string
	push       bool
	reportPath string
	vexPath    string
	cmd        *exec.Cmd // Current command being built
}

type PatchParamsConstraint interface {
	types.ReportBasedPatchParams |
		types.PlatformSelectivePatchParams |
		types.ComprehensivePatchParams |
		types.PatchParams |
		types.BasePatchParams
}

// NOTE: use generic for param types to assist the agent with populating the correct values.
func New[T PatchParamsConstraint](params T, dryRun bool) *CLI {
	var image, tag, reportPath string
	var platforms []string
	var push bool

	// Extract common fields using type switch
	switch p := any(params).(type) {
	case types.ReportBasedPatchParams:
		image, tag, push, reportPath = p.Image, p.Tag, p.Push, p.ReportPath
	case types.PlatformSelectivePatchParams:
		image, tag, push, platforms = p.Image, p.Tag, p.Push, p.Platform
	case types.ComprehensivePatchParams:
		image, tag, push = p.Image, p.Tag, p.Push
	case types.PatchParams:
		image, tag, push, platforms = p.Image, p.Tag, p.Push, p.Platform
	case types.BasePatchParams:
		image, tag, push = p.Image, p.Tag, p.Push
	}

	return &CLI{
		copaPath:   "copa",
		dryRun:     dryRun,
		image:      image,
		tag:        tag,
		platforms:  platforms,
		push:       push,
		reportPath: reportPath,
	}
}

func (c *CLI) Build() *CLI {
	args := []string{"patch"}
	args = append(args, "--image", c.image)

	if c.tag != "" {
		args = append(args, "--tag", c.tag)
	}

	if c.push {
		args = append(args, "--push")
	}

	c.cmd = exec.Command(c.copaPath, args...)
	return c
}

func (c *CLI) BuildWithPlatforms() *CLI {
	c = c.Build()

	if len(c.platforms) > 0 {
		supportedPlatforms := FilterSupportedPlatforms(c.platforms)
		if len(supportedPlatforms) > 0 {
			c.cmd.Args = append(c.cmd.Args, "--platform", strings.Join(supportedPlatforms, ","))
		}
	}

	return c
}

func (c *CLI) BuildWithReport() *CLI {
	c = c.Build()

	if c.reportPath != "" {
		c.cmd.Args = append(c.cmd.Args, "--report", c.reportPath)
		c.vexPath = filepath.Join(os.TempDir(), defaultVexFile)
		c.cmd.Args = append(c.cmd.Args, "--output", c.vexPath)
	}

	return c
}

func (c *CLI) Run(ctx context.Context) error {

	remotePatch, err := docker.SetupRegistryAuthFromEnv()
	if err != nil {
		return fmt.Errorf("failed to autenticate to registry: %w", err)
	}
	if remotePatch {
		c.push = true
	}

	if c.cmd == nil {
		return fmt.Errorf("no command built - call a Build method first")
	}

	if c.dryRun {
		fmt.Printf("[DRY RUN] %s %s\n", c.cmd.Path, strings.Join(c.cmd.Args[1:], " "))
		return nil
	}

	fmt.Printf("Executing: %s %s\n", c.cmd.Path, strings.Join(c.cmd.Args[1:], " "))

	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	return c.cmd.Run()
}

func (c *CLI) RunOutputVex(ctx context.Context) (string, error) {
	if err := c.Run(ctx); err != nil {
		return "", err
	}
	return c.vexPath, nil
}

// IsPlatformSupported checks if the given platform is supported by Copa for patching
func IsPlatformSupported(platform string) bool {
	for _, supported := range CopaSupportedPlatforms {
		if platform == supported {
			return true
		}
		// Handle arm64 variants - Copa supports "linux/arm64" which covers "linux/arm64/v8"
		if supported == "linux/arm64" && (platform == "linux/arm64/v8" || platform == "linux/arm64") {
			return true
		}
	}
	return false
}

// FilterSupportedPlatforms returns only the platforms that Copa can patch from the given list
func FilterSupportedPlatforms(platforms []string) []string {
	var supported []string
	for _, platform := range platforms {
		if IsPlatformSupported(platform) {
			supported = append(supported, platform)
		}
	}
	return supported
}
