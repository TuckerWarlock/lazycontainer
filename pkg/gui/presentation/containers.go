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
		getDisplayCPUInfo(container),
		utils.ColoredString(container.GetPorts(), color.FgYellow),
		utils.ColoredString(container.GetImage(), color.FgMagenta),
	}
}

// getContainerDisplayStatus returns the colored status indicator
func getContainerDisplayStatus(c *commands.Container) string {
	statusMap := map[string]struct {
		symbol string
		color  color.Attribute
	}{
		"running": {"▶", color.FgGreen},
		"stopped": {"⨯", color.FgRed},
		"exited":  {"⨯", color.FgYellow},
		"paused":  {"◫", color.FgYellow},
		"created": {"+", color.FgCyan},
	}

	if status, ok := statusMap[c.GetStatus()]; ok {
		return utils.ColoredString(status.symbol, status.color)
	}

	return utils.ColoredString("?", color.FgWhite)
}

// getDisplayCPUInfo returns CPU allocation info
func getDisplayCPUInfo(c *commands.Container) string {
	if !c.IsRunning() {
		return ""
	}
	return c.GetMemoryHuman()
}

// getContainerColor returns the color attribute for a container based on its state
func getContainerColor(c *commands.Container) color.Attribute {
	switch c.GetStatus() {
	case "running":
		return color.FgGreen
	case "stopped", "exited":
		return color.FgYellow
	case "paused":
		return color.FgYellow
	case "created":
		return color.FgCyan
	default:
		return color.FgWhite
	}
}
