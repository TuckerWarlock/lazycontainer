package presentation

import (
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

// GetContainerDisplayStrings returns the display strings for a container using the default theme
func GetContainerDisplayStrings(container *commands.Container) []string {
	return GetContainerDisplayStringsWithTheme(container, ThemeDefault)
}

// GetContainerDisplayStringsWithTheme returns the display strings for a container with a specific theme
func GetContainerDisplayStringsWithTheme(container *commands.Container, theme *Theme) []string {
	return []string{
		getContainerDisplayStatusWithTheme(container, theme),
		container.GetName(),
		getContainerIPOrUptime(container),
		utils.ColoredString(container.GetPorts(), theme.GetPortsColor()),
		utils.ColoredString(getShortImage(container.GetImage()), theme.GetImageColor()),
	}
}

// getContainerDisplayStatus returns the colored status indicator
func getContainerDisplayStatus(c *commands.Container) string {
	return getContainerDisplayStatusWithTheme(c, ThemeDefault)
}

// getContainerDisplayStatusWithTheme returns the colored status indicator with a specific theme
func getContainerDisplayStatusWithTheme(c *commands.Container, theme *Theme) string {
	symbol := c.GetStatusSymbol()
	return utils.ColoredString(symbol, theme.GetContainerStatusColor(c.GetStatus()))
}

// getContainerIPOrUptime returns IP address for running containers, or status for stopped
func getContainerIPOrUptime(c *commands.Container) string {
	if c.IsRunning() {
		ip := c.GetIPAddress()
		if ip != "" {
			return ip
		}
		return c.GetUptime()
	}
	return c.GetStatus()
}

// getShortImage returns a shortened image name
func getShortImage(image string) string {
	// Remove common prefixes
	prefixes := []string{
		"docker.io/library/",
		"docker.io/",
		"ghcr.io/",
		"quay.io/",
	}
	for _, prefix := range prefixes {
		if len(image) > len(prefix) && image[:len(prefix)] == prefix {
			image = image[len(prefix):]
			break
		}
	}
	// Truncate if too long
	if len(image) > 25 {
		return image[:22] + "..."
	}
	return image
}
