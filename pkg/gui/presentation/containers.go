package presentation

import (
	"github.com/fatih/color"
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

// GetContainerDisplayStrings returns the display strings for a container
func GetContainerDisplayStrings(container *commands.Container) []string {
	return []string{
		getContainerDisplayStatus(container),
		container.GetName(),
		getContainerIPOrUptime(container),
		utils.ColoredString(container.GetPorts(), color.FgYellow),
		utils.ColoredString(getShortImage(container.GetImage()), color.FgMagenta),
	}
}

// getContainerDisplayStatus returns the colored status indicator
func getContainerDisplayStatus(c *commands.Container) string {
	symbol := c.GetStatusSymbol()
	return utils.ColoredString(symbol, getContainerColor(c))
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

// getContainerColor returns the color attribute for a container based on its state
func getContainerColor(c *commands.Container) color.Attribute {
	switch c.GetStatus() {
	case "running":
		return color.FgGreen
	case "stopped", "exited":
		return color.FgRed
	case "paused":
		return color.FgYellow
	case "created":
		return color.FgCyan
	default:
		return color.FgWhite
	}
}
