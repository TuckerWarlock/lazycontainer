package presentation

import (
	"github.com/warl0ck/lazycontainer/pkg/commands"
)

// GetImageDisplayStrings returns the display strings for an image
func GetImageDisplayStrings(image *commands.Image) []string {
	return []string{
		getShortImageRef(image.Reference),
		image.GetDigestShort(),
		image.GetSizeHuman(),
	}
}

// getShortImageRef returns a shortened image reference
func getShortImageRef(ref string) string {
	// Remove common prefixes
	prefixes := []string{
		"docker.io/library/",
		"docker.io/",
		"ghcr.io/",
		"quay.io/",
	}
	for _, prefix := range prefixes {
		if len(ref) > len(prefix) && ref[:len(prefix)] == prefix {
			ref = ref[len(prefix):]
			break
		}
	}
	return ref
}
