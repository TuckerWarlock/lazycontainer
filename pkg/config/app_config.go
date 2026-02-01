package config

import (
	"os"
	"path/filepath"
)

// AppConfig holds application configuration
type AppConfig struct {
	Name        string
	Version     string
	Commit      string
	Date        string
	BuildSource string
	Debug       bool
	ConfigDir   string
	UserConfig  *UserConfig
}

// UserConfig holds user-customizable settings
type UserConfig struct {
	Gui           GuiConfig           `yaml:"gui"`
	Logs          LogsConfig          `yaml:"logs"`
	Stats         StatsConfig         `yaml:"stats"`
	Confirmations ConfirmationsConfig `yaml:"confirmations"`
}

// GuiConfig holds GUI-related settings
type GuiConfig struct {
	ScrollHeight      int    `yaml:"scrollHeight"`
	ScrollPastBottom  bool   `yaml:"scrollPastBottom"`
	MouseEvents       bool   `yaml:"mouseEvents"`
	Theme             Theme  `yaml:"theme"`
	ShowAllContainers bool   `yaml:"showAllContainers"`
	ReturnImmediately bool   `yaml:"returnImmediately"`
	WrapMainPanel     bool   `yaml:"wrapMainPanel"`
}

// Theme holds color theme settings
type Theme struct {
	ActiveBorderColor   string `yaml:"activeBorderColor"`
	InactiveBorderColor string `yaml:"inactiveBorderColor"`
	SelectedLineBgColor string `yaml:"selectedLineBgColor"`
}

// LogsConfig holds log viewing settings
type LogsConfig struct {
	Timestamps bool `yaml:"timestamps"`
	Since      string `yaml:"since"`
	Tail       int    `yaml:"tail"`
}

// StatsConfig holds stats display settings
type StatsConfig struct {
	ShowAllContainers bool `yaml:"showAllContainers"`
}

// ConfirmationsConfig holds confirmation dialog settings
type ConfirmationsConfig struct {
	Remove bool `yaml:"remove"`
	Stop   bool `yaml:"stop"`
	Kill   bool `yaml:"kill"`
}

// NewAppConfig creates a new AppConfig
func NewAppConfig(name, version, commit, date, buildSource string, debug bool) (*AppConfig, error) {
	configDir, err := getConfigDir(name)
	if err != nil {
		return nil, err
	}

	userConfig := GetDefaultConfig()

	return &AppConfig{
		Name:        name,
		Version:     version,
		Commit:      commit,
		Date:        date,
		BuildSource: buildSource,
		Debug:       debug,
		ConfigDir:   configDir,
		UserConfig:  userConfig,
	}, nil
}

// GetDefaultConfig returns default configuration
func GetDefaultConfig() *UserConfig {
	return &UserConfig{
		Gui: GuiConfig{
			ScrollHeight:      2,
			ScrollPastBottom:  true,
			MouseEvents:       true,
			ShowAllContainers: true,
			ReturnImmediately: false,
			WrapMainPanel:     true,
			Theme: Theme{
				ActiveBorderColor:   "green",
				InactiveBorderColor: "default",
				SelectedLineBgColor: "blue",
			},
		},
		Logs: LogsConfig{
			Timestamps: false,
			Since:      "",
			Tail:       100,
		},
		Stats: StatsConfig{
			ShowAllContainers: false,
		},
		Confirmations: ConfirmationsConfig{
			Remove: true,
			Stop:   false,
			Kill:   false,
		},
	}
}

// getConfigDir returns the configuration directory path
func getConfigDir(appName string) (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(home, ".config")
	}

	configDir := filepath.Join(configHome, appName)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetLogPath returns the path to the log file
func (c *AppConfig) GetLogPath() string {
	return filepath.Join(c.ConfigDir, "development.log")
}
