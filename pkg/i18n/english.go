package i18n

// TranslationSet is a set of localised strings for a given language
type TranslationSet struct {
	NotEnoughSpace   string
	ProjectTitle     string
	MainTitle        string
	GlobalTitle      string
	Navigate         string
	Menu             string
	MenuTitle        string
	Execute          string
	Scroll           string
	Close            string
	Quit             string
	ErrorTitle       string
	OpenConfig       string
	EditConfig       string
	ConfirmQuit      string
	ErrorOccurred    string
	ConnectionFailed string

	CannotKillChildError string

	Donate             string
	Cancel             string
	CustomCommandTitle string
	BulkCommandTitle   string
	Remove             string
	HideStopped        string
	ForceRemove        string
	Confirm            string
	Return             string
	FocusMain          string
	LcFilter           string
	StopContainer      string
	RestartingStatus   string
	StartingStatus     string
	StoppingStatus     string
	PausingStatus      string
	RemovingStatus     string
	PruningStatus      string

	RunningCustomCommandStatus string
	RunningBulkCommandStatus   string

	Stop            string
	Pause           string
	Restart         string
	Start           string
	PreviousContext string
	NextContext     string
	Attach          string
	ViewLogs        string

	ContainersTitle string
	TopTitle        string
	ImagesTitle     string
	VolumesTitle    string
	NetworksTitle   string
	NoContainers    string
	NoContainer     string
	NoImages        string
	NoVolumes       string
	NoNetworks      string
	RemoveImage     string
	RemoveVolume    string
	RemoveNetwork   string
	RemoveWithForce string
	PruneImages     string
	PruneContainers string
	PruneVolumes    string
	PruneNetworks   string

	ConfirmPruneContainers  string
	ConfirmStopContainers   string
	ConfirmRemoveContainers string
	ConfirmPruneImages      string
	ConfirmPruneVolumes     string
	ConfirmPruneNetworks    string

	StopAllContainers    string
	RemoveAllContainers  string
	ViewRestartOptions   string
	ExecShell            string
	RunCustomCommand     string
	ViewBulkCommands     string
	FilterList           string
	OpenInBrowser        string
	SortContainersByState string

	LogsTitle         string
	ConfigTitle       string
	EnvTitle          string
	StatsTitle        string
	CreditsTitle      string
	NothingToDisplay  string

	No  string
	Yes string

	LcNextScreenMode string
	LcPrevScreenMode string
	FilterPrompt     string

	FocusContainers string
	FocusImages     string
	FocusVolumes    string
	FocusNetworks   string
}

func englishSet() TranslationSet {
	return TranslationSet{
		PruningStatus:              "pruning",
		RemovingStatus:             "removing",
		RestartingStatus:           "restarting",
		StartingStatus:             "starting",
		StoppingStatus:             "stopping",
		PausingStatus:              "pausing",
		RunningCustomCommandStatus: "running custom command",
		RunningBulkCommandStatus:   "running bulk command",

		ErrorOccurred:    "An error occurred! Please create an issue at https://github.com/warl0ck/lazycontainer/issues",
		ConnectionFailed: "Connection to container service failed. Make sure the container service is running with: container system start",

		CannotKillChildError: "Waited three seconds for child process to stop. There may be an orphan process that continues to run on your system.",

		Donate:  "Donate",
		Confirm: "Confirm",

		Return:       "return",
		FocusMain:    "focus main panel",
		LcFilter:     "filter list",
		Navigate:     "navigate",
		Execute:      "execute",
		Close:        "close",
		Quit:         "quit",
		Menu:         "menu",
		MenuTitle:    "Menu",
		Scroll:       "scroll",
		OpenConfig:   "open lazycontainer config",
		EditConfig:   "edit lazycontainer config",
		Cancel:       "cancel",
		Remove:       "remove",
		HideStopped:  "hide/show stopped containers",
		ForceRemove:  "force remove",
		Stop:         "stop",
		Pause:        "pause",
		Restart:      "restart",
		Start:        "start",
		PreviousContext: "previous tab",
		NextContext:     "next tab",
		Attach:          "attach",
		ViewLogs:        "view logs",
		RemoveImage:     "remove image",
		RemoveVolume:    "remove volume",
		RemoveNetwork:   "remove network",
		RemoveWithForce: "remove (forced)",
		PruneContainers: "prune exited containers",
		PruneVolumes:    "prune unused volumes",
		PruneNetworks:   "prune unused networks",
		PruneImages:     "prune unused images",
		StopAllContainers:   "stop all containers",
		RemoveAllContainers: "remove all containers (forced)",
		ViewRestartOptions:  "view restart options",
		ExecShell:           "exec shell",
		RunCustomCommand:    "run predefined custom command",
		ViewBulkCommands:    "view bulk commands",
		FilterList:          "filter list",
		OpenInBrowser:       "open in browser (first port is http)",
		SortContainersByState: "sort containers by state",

		GlobalTitle:     "Global",
		MainTitle:       "Main",
		ProjectTitle:    "Project",
		ContainersTitle: "Containers",
		ImagesTitle:     "Images",
		VolumesTitle:    "Volumes",
		NetworksTitle:   "Networks",
		CustomCommandTitle: "Custom Command:",
		BulkCommandTitle:   "Bulk Command:",
		ErrorTitle:         "Error",
		LogsTitle:          "Logs",
		ConfigTitle:        "Config",
		EnvTitle:           "Env",
		TopTitle:           "Top",
		StatsTitle:         "Stats",
		CreditsTitle:       "About",
		NothingToDisplay:   "Nothing to display",

		NoContainers: "No containers",
		NoContainer:  "No container",
		NoImages:     "No images",
		NoVolumes:    "No volumes",
		NoNetworks:   "No networks",

		ConfirmQuit:            "Are you sure you want to quit?",
		NotEnoughSpace:         "Not enough space to render panels",
		ConfirmPruneImages:     "Are you sure you want to prune all unused images?",
		ConfirmPruneContainers: "Are you sure you want to prune all stopped containers?",
		ConfirmStopContainers:  "Are you sure you want to stop all containers?",
		ConfirmRemoveContainers: "Are you sure you want to remove all containers?",
		ConfirmPruneVolumes:    "Are you sure you want to prune all unused volumes?",
		ConfirmPruneNetworks:   "Are you sure you want to prune all unused networks?",
		StopContainer:          "Are you sure you want to stop this container?",

		No:  "no",
		Yes: "yes",

		LcNextScreenMode: "next screen mode (normal/half/fullscreen)",
		LcPrevScreenMode: "prev screen mode",
		FilterPrompt:     "filter",

		FocusContainers: "focus containers panel",
		FocusImages:     "focus images panel",
		FocusVolumes:    "focus volumes panel",
		FocusNetworks:   "focus networks panel",
	}
}
