package copa

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

// ImageInfo contains information about an image's platform support and availability
type ImageInfo struct {
	IsMultiPlatform bool
	IsLocal         bool
	Platform        []string // Available platforms (e.g., ["linux/amd64", "linux/arm64"])
}
