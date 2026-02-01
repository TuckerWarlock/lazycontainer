package presentation

import (
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

// GetVolumeDisplayStrings returns the display strings for a volume
func GetVolumeDisplayStrings(volume *commands.Volume) []string {
	return []string{
		volume.Name,
		volume.Format,
		utils.FormatBinaryBytes(int(volume.SizeInBytes)),
	}
}
