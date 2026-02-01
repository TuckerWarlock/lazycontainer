package presentation

import (
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

// GetImageDisplayStrings returns the display strings for an image
func GetImageDisplayStrings(image *commands.Image) []string {
	return []string{
		image.Reference,
		utils.SafeTruncate(image.GetDigestShort(), 12),
		utils.FormatBinaryBytes(int(image.Size)),
	}
}
