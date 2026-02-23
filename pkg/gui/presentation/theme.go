package presentation

import (
	"github.com/fatih/color"
)

// Theme defines colors for the UI
type Theme struct {
	Name string

	// Container status colors
	ContainerRunning color.Attribute
	ContainerStopped color.Attribute
	ContainerPaused  color.Attribute
	ContainerCreated color.Attribute
	ContainerDefault color.Attribute

	// Secondary colors for details
	Ports color.Attribute
	Image color.Attribute
}

// ThemeDefault is the default dark theme
var ThemeDefault = &Theme{
	Name:             "default",
	ContainerRunning: color.FgGreen,
	ContainerStopped: color.FgRed,
	ContainerPaused:  color.FgYellow,
	ContainerCreated: color.FgCyan,
	ContainerDefault: color.FgWhite,
	Ports:            color.FgYellow,
	Image:            color.FgMagenta,
}

// ThemeLight is a light-friendly theme with muted colors
var ThemeLight = &Theme{
	Name:             "light",
	ContainerRunning: color.FgGreen,
	ContainerStopped: color.FgRed,
	ContainerPaused:  color.FgYellow,
	ContainerCreated: color.FgBlue,
	ContainerDefault: color.FgBlack,
	Ports:            color.FgYellow,
	Image:            color.FgCyan,
}

// ThemeDark is an enhanced dark theme with high contrast
var ThemeDark = &Theme{
	Name:             "dark",
	ContainerRunning: color.FgGreen,
	ContainerStopped: color.FgRed,
	ContainerPaused:  color.FgYellow,
	ContainerCreated: color.FgMagenta,
	ContainerDefault: color.FgWhite,
	Ports:            color.FgCyan,
	Image:            color.FgYellow,
}

// GetContainerStatusColor returns the color for a container status from the given theme
func (t *Theme) GetContainerStatusColor(status string) color.Attribute {
	switch status {
	case "running":
		return t.ContainerRunning
	case "stopped", "exited":
		return t.ContainerStopped
	case "paused":
		return t.ContainerPaused
	case "created":
		return t.ContainerCreated
	default:
		return t.ContainerDefault
	}
}

// GetPortsColor returns the color for ports from the given theme
func (t *Theme) GetPortsColor() color.Attribute {
	return t.Ports
}

// GetImageColor returns the color for image names from the given theme
func (t *Theme) GetImageColor() color.Attribute {
	return t.Image
}

// AllThemes returns all available themes
func AllThemes() []*Theme {
	return []*Theme{ThemeDefault, ThemeLight, ThemeDark}
}

// GetThemeByName returns a theme by name, or ThemeDefault if not found
func GetThemeByName(name string) *Theme {
	for _, theme := range AllThemes() {
		if theme.Name == name {
			return theme
		}
	}
	return ThemeDefault
}
