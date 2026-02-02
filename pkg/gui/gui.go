package gui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/config"
	"github.com/warl0ck/lazycontainer/pkg/gui/panels"
	"github.com/warl0ck/lazycontainer/pkg/gui/presentation"
	"github.com/warl0ck/lazycontainer/pkg/gui/types"
	"github.com/warl0ck/lazycontainer/pkg/i18n"
	"github.com/warl0ck/lazycontainer/pkg/tasks"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

// WindowMaximisation determines how much space your selected window takes up
type WindowMaximisation int

const (
	SCREEN_NORMAL WindowMaximisation = iota
	SCREEN_HALF
	SCREEN_FULL
)

// Panels holds all the side list panels
type Panels struct {
	Containers *panels.SideListPanel[*commands.Container]
	Images     *panels.SideListPanel[*commands.Image]
	Volumes    *panels.SideListPanel[*commands.Volume]
	Networks   *panels.SideListPanel[*commands.Network]
	Menu       *panels.SideListPanel[*types.MenuItem]
}

type mainPanelState struct {
	ObjectKey string
}

type panelStates struct {
	Main *mainPanelState
}

type filterState struct {
	active bool
	panel  panels.ISideListPanel
	needle string
}

type guiState struct {
	ViewStack            []string
	Panels               *panelStates
	ShowExitedContainers bool
	ScreenMode           WindowMaximisation
	Filter               filterState
}

// Gui wraps the gocui Gui object
type Gui struct {
	g                *gocui.Gui
	Log              *logrus.Entry
	ContainerCommand *commands.ContainerCommand
	OSCommand        *commands.OSCommand
	Config           *config.AppConfig
	Tr               *i18n.TranslationSet

	Views       Views
	Panels      Panels
	State       guiState
	taskManager *tasks.TaskManager

	statusMessage string
	lastRefresh   time.Time
}

// NewGui creates a new Gui instance
func NewGui(log *logrus.Entry, containerCmd *commands.ContainerCommand, osCmd *commands.OSCommand, config *config.AppConfig) (*Gui, error) {
	tr := i18n.GetTranslationSet()

	initialState := guiState{
		Panels: &panelStates{
			Main: &mainPanelState{
				ObjectKey: "",
			},
		},
		ViewStack:            []string{},
		ShowExitedContainers: true,
		ScreenMode:           SCREEN_NORMAL,
	}

	gui := &Gui{
		Log:              log,
		ContainerCommand: containerCmd,
		OSCommand:        osCmd,
		Config:           config,
		Tr:               tr,
		State:            initialState,
		taskManager:      tasks.NewTaskManager(log, tr),
	}

	return gui, nil
}

// Run starts the GUI main loop
func (gui *Gui) Run() error {
	defer gui.taskManager.Close()

	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode: gocui.OutputTrue,
	})
	if err != nil {
		return fmt.Errorf("failed to create GUI: %w", err)
	}
	defer g.Close()

	gui.g = g
	g.Mouse = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.SelFrameColor = gocui.ColorGreen

	g.SetManagerFunc(gui.layout)

	if err := gui.initializeViews(); err != nil {
		return fmt.Errorf("failed to initialize views: %w", err)
	}

	gui.setPanels()

	if err := gui.keybindings(g); err != nil {
		return fmt.Errorf("failed to set keybindings: %w", err)
	}

	if gui.g.CurrentView() == nil {
		viewName := gui.initiallyFocusedViewName()
		view, err := gui.g.View(viewName)
		if err != nil {
			return fmt.Errorf("failed to get initial view: %w", err)
		}
		if err := gui.switchFocus(view); err != nil {
			return fmt.Errorf("failed to set initial focus: %w", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go gui.refresh()
	go gui.backgroundRefresh(ctx)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (gui *Gui) initiallyFocusedViewName() string {
	return "containers"
}

func (gui *Gui) setPanels() {
	gui.Panels = Panels{
		Containers: gui.getContainersPanel(),
		Images:     gui.getImagesPanel(),
		Volumes:    gui.getVolumesPanel(),
		Networks:   gui.getNetworksPanel(),
		Menu:       gui.getMenuPanel(),
	}
}

func (gui *Gui) getContainersPanel() *panels.SideListPanel[*commands.Container] {
	return &panels.SideListPanel[*commands.Container]{
		ContextState: &panels.ContextState[*commands.Container]{
			GetMainTabs: func() []panels.MainTab[*commands.Container] {
				return []panels.MainTab[*commands.Container]{
					{
						Key:    "logs",
						Title:  gui.Tr.LogsTitle,
						Render: gui.renderContainerLogsTask,
					},
					{
						Key:    "config",
						Title:  gui.Tr.ConfigTitle,
						Render: gui.renderContainerConfigTask,
					},
					{
						Key:    "stats",
						Title:  gui.Tr.StatsTitle,
						Render: gui.renderContainerStatsTask,
					},
				}
			},
			GetItemContextCacheKey: func(c *commands.Container) string {
				return "containers-" + c.GetID() + "-" + c.GetStatus()
			},
		},
		ListPanel: panels.ListPanel[*commands.Container]{
			List: panels.NewFilteredList[*commands.Container](),
			View: gui.Views.Containers,
		},
		NoItemsMessage: gui.Tr.NoContainers,
		Gui:            gui,
		GetTableCells:  presentation.GetContainerDisplayStrings,
	}
}

func (gui *Gui) getImagesPanel() *panels.SideListPanel[*commands.Image] {
	return &panels.SideListPanel[*commands.Image]{
		ContextState: &panels.ContextState[*commands.Image]{
			GetMainTabs: func() []panels.MainTab[*commands.Image] {
				return []panels.MainTab[*commands.Image]{
					{
						Key:    "config",
						Title:  gui.Tr.ConfigTitle,
						Render: gui.renderImageConfigTask,
					},
				}
			},
			GetItemContextCacheKey: func(i *commands.Image) string {
				return "images-" + i.GetDigest()
			},
		},
		ListPanel: panels.ListPanel[*commands.Image]{
			List: panels.NewFilteredList[*commands.Image](),
			View: gui.Views.Images,
		},
		NoItemsMessage: gui.Tr.NoImages,
		Gui:            gui,
		GetTableCells:  presentation.GetImageDisplayStrings,
	}
}

func (gui *Gui) getVolumesPanel() *panels.SideListPanel[*commands.Volume] {
	return &panels.SideListPanel[*commands.Volume]{
		ContextState: &panels.ContextState[*commands.Volume]{
			GetMainTabs: func() []panels.MainTab[*commands.Volume] {
				return []panels.MainTab[*commands.Volume]{
					{
						Key:    "config",
						Title:  gui.Tr.ConfigTitle,
						Render: gui.renderVolumeConfigTask,
					},
				}
			},
			GetItemContextCacheKey: func(v *commands.Volume) string {
				return "volumes-" + v.Name
			},
		},
		ListPanel: panels.ListPanel[*commands.Volume]{
			List: panels.NewFilteredList[*commands.Volume](),
			View: gui.Views.Volumes,
		},
		NoItemsMessage: gui.Tr.NoVolumes,
		Gui:            gui,
		GetTableCells:  presentation.GetVolumeDisplayStrings,
	}
}

func (gui *Gui) getNetworksPanel() *panels.SideListPanel[*commands.Network] {
	return &panels.SideListPanel[*commands.Network]{
		ContextState: &panels.ContextState[*commands.Network]{
			GetMainTabs: func() []panels.MainTab[*commands.Network] {
				return []panels.MainTab[*commands.Network]{
					{
						Key:    "config",
						Title:  gui.Tr.ConfigTitle,
						Render: gui.renderNetworkConfigTask,
					},
				}
			},
			GetItemContextCacheKey: func(n *commands.Network) string {
				return "networks-" + n.ID
			},
		},
		ListPanel: panels.ListPanel[*commands.Network]{
			List: panels.NewFilteredList[*commands.Network](),
			View: gui.Views.Networks,
		},
		NoItemsMessage: gui.Tr.NoNetworks,
		Gui:            gui,
		GetTableCells:  presentation.GetNetworkDisplayStrings,
	}
}

func (gui *Gui) getMenuPanel() *panels.SideListPanel[*types.MenuItem] {
	return &panels.SideListPanel[*types.MenuItem]{
		ListPanel: panels.ListPanel[*types.MenuItem]{
			List: panels.NewFilteredList[*types.MenuItem](),
			View: gui.Views.Menu,
		},
		Gui:           gui,
		DisableFilter: true,
		GetTableCells: func(item *types.MenuItem) []string {
			return item.LabelColumns
		},
		OnClick: func(item *types.MenuItem) error {
			if item.OnPress != nil {
				return item.OnPress()
			}
			return nil
		},
	}
}

// render task functions for each panel type
func (gui *Gui) renderContainerLogsTask(c *commands.Container) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if c == nil {
			return ""
		}
		logs, err := gui.ContainerCommand.GetContainerLogs(c.GetID(), 100)
		if err != nil {
			return fmt.Sprintf("Error fetching logs: %v", err)
		}
		return logs
	})
}

func (gui *Gui) renderContainerConfigTask(c *commands.Container) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if c == nil {
			return ""
		}
		return gui.formatContainerConfig(c)
	})
}

func (gui *Gui) renderContainerStatsTask(c *commands.Container) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if c == nil {
			return ""
		}
		return gui.formatContainerStats(c)
	})
}

func (gui *Gui) renderImageConfigTask(i *commands.Image) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if i == nil {
			return ""
		}
		return gui.formatImageConfig(i)
	})
}

func (gui *Gui) renderVolumeConfigTask(v *commands.Volume) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if v == nil {
			return ""
		}
		return gui.formatVolumeConfig(v)
	})
}

func (gui *Gui) renderNetworkConfigTask(n *commands.Network) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if n == nil {
			return ""
		}
		return gui.formatNetworkConfig(n)
	})
}

func (gui *Gui) formatContainerConfig(c *commands.Container) string {
	var sb strings.Builder
	sb.WriteString(utils.ColoredString("Container Configuration\n", color.FgCyan))
	sb.WriteString(strings.Repeat("─", 40) + "\n\n")

	sb.WriteString(utils.FormatMapItem(0, "Name", c.GetName()))
	sb.WriteString(utils.FormatMapItem(0, "Status", c.GetStatus()))
	sb.WriteString(utils.FormatMapItem(0, "Image", c.GetImage()))
	sb.WriteString(utils.FormatMapItem(0, "CPUs", c.GetCPUs()))
	sb.WriteString(utils.FormatMapItem(0, "Memory", c.GetMemoryHuman()))

	if ports := c.GetPorts(); ports != "" {
		sb.WriteString(utils.FormatMapItem(0, "Ports", ports))
	}

	if c.IsRunning() {
		sb.WriteString(utils.FormatMapItem(0, "Uptime", c.GetUptime()))
		sb.WriteString(utils.FormatMapItem(0, "Started", c.GetStartedAt().Format("2006-01-02 15:04:05")))
	}

	return sb.String()
}

func (gui *Gui) formatContainerStats(c *commands.Container) string {
	var sb strings.Builder
	sb.WriteString(utils.ColoredString("Container Stats\n", color.FgCyan))
	sb.WriteString(strings.Repeat("─", 40) + "\n\n")

	if !c.IsRunning() {
		sb.WriteString("Container is not running\n")
		return sb.String()
	}

	sb.WriteString(utils.FormatMapItem(0, "Status", c.GetStatus()))
	sb.WriteString(utils.FormatMapItem(0, "CPUs Allocated", c.GetCPUs()))
	sb.WriteString(utils.FormatMapItem(0, "Memory Allocated", c.GetMemoryHuman()))
	sb.WriteString(utils.FormatMapItem(0, "Uptime", c.GetUptime()))

	return sb.String()
}

func (gui *Gui) formatImageConfig(i *commands.Image) string {
	var sb strings.Builder
	sb.WriteString(utils.ColoredString("Image Details\n", color.FgCyan))
	sb.WriteString(strings.Repeat("─", 40) + "\n\n")

	sb.WriteString(utils.FormatMapItem(0, "Reference", i.Reference))
	sb.WriteString(utils.FormatMapItem(0, "Digest", i.GetDigest()))
	sb.WriteString(utils.FormatMapItem(0, "Size", i.GetSizeHuman()))

	return sb.String()
}

func (gui *Gui) formatVolumeConfig(v *commands.Volume) string {
	var sb strings.Builder
	sb.WriteString(utils.ColoredString("Volume Details\n", color.FgCyan))
	sb.WriteString(strings.Repeat("─", 40) + "\n\n")

	sb.WriteString(utils.FormatMapItem(0, "Name", v.Name))
	sb.WriteString(utils.FormatMapItem(0, "Format", v.Format))
	sb.WriteString(utils.FormatMapItem(0, "Size", v.GetSizeHuman()))

	if v.Source != "" {
		sb.WriteString(utils.FormatMapItem(0, "Source", v.Source))
	}

	return sb.String()
}

func (gui *Gui) formatNetworkConfig(n *commands.Network) string {
	var sb strings.Builder
	sb.WriteString(utils.ColoredString("Network Details\n", color.FgCyan))
	sb.WriteString(strings.Repeat("─", 40) + "\n\n")

	sb.WriteString(utils.FormatMapItem(0, "ID", n.ID))
	sb.WriteString(utils.FormatMapItem(0, "Mode", n.Config.Mode))

	if n.Status.IPv4Subnet != "" {
		sb.WriteString(utils.FormatMapItem(0, "IPv4 Subnet", n.Status.IPv4Subnet))
	}
	if n.Status.IPv4Gateway != "" {
		sb.WriteString(utils.FormatMapItem(0, "IPv4 Gateway", n.Status.IPv4Gateway))
	}
	if n.Status.IPv6Subnet != "" {
		sb.WriteString(utils.FormatMapItem(0, "IPv6 Subnet", n.Status.IPv6Subnet))
	}

	return sb.String()
}

func (gui *Gui) backgroundRefresh(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			gui.refresh()
		}
	}
}

func (gui *Gui) refresh() {
	gui.lastRefresh = time.Now()

	go func() {
		if err := gui.refreshContainers(); err != nil {
			gui.Log.Error(err)
		}
	}()
	go func() {
		if err := gui.refreshImages(); err != nil {
			gui.Log.Error(err)
		}
	}()
	go func() {
		if err := gui.refreshVolumes(); err != nil {
			gui.Log.Error(err)
		}
	}()
	go func() {
		if err := gui.refreshNetworks(); err != nil {
			gui.Log.Error(err)
		}
	}()
}

func (gui *Gui) refreshContainers() error {
	containers, err := gui.ContainerCommand.ListContainers(gui.State.ShowExitedContainers)
	if err != nil {
		return err
	}

	// Convert to pointer slice
	containerPtrs := make([]*commands.Container, len(containers))
	for i := range containers {
		containerPtrs[i] = &containers[i]
	}

	gui.Panels.Containers.SetItems(containerPtrs)
	return gui.Panels.Containers.RerenderList()
}

func (gui *Gui) refreshImages() error {
	images, err := gui.ContainerCommand.ListImages()
	if err != nil {
		return err
	}

	imagePtrs := make([]*commands.Image, len(images))
	for i := range images {
		imagePtrs[i] = &images[i]
	}

	gui.Panels.Images.SetItems(imagePtrs)
	return gui.Panels.Images.RerenderList()
}

func (gui *Gui) refreshVolumes() error {
	volumes, err := gui.ContainerCommand.ListVolumes()
	if err != nil {
		return err
	}

	volumePtrs := make([]*commands.Volume, len(volumes))
	for i := range volumes {
		volumePtrs[i] = &volumes[i]
	}

	gui.Panels.Volumes.SetItems(volumePtrs)
	return gui.Panels.Volumes.RerenderList()
}

func (gui *Gui) refreshNetworks() error {
	networks, err := gui.ContainerCommand.ListNetworks()
	if err != nil {
		return err
	}

	networkPtrs := make([]*commands.Network, len(networks))
	for i := range networks {
		networkPtrs[i] = &networks[i]
	}

	gui.Panels.Networks.SetItems(networkPtrs)
	return gui.Panels.Networks.RerenderList()
}

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Minimum terminal size check
	if maxX < 60 || maxY < 15 {
		gui.showLimitView(g, maxX, maxY)
		return nil
	}
	gui.hideLimitView()

	// Calculate dimensions
	sidePanelWidth := maxX * 30 / 100
	if sidePanelWidth < 30 {
		sidePanelWidth = 30
	}
	if sidePanelWidth > 60 {
		sidePanelWidth = 60
	}

	// Apply screen mode adjustments
	switch gui.State.ScreenMode {
	case SCREEN_HALF:
		sidePanelWidth = maxX * 15 / 100
		if sidePanelWidth < 20 {
			sidePanelWidth = 20
		}
	case SCREEN_FULL:
		sidePanelWidth = 0
	}

	bottomHeight := 2
	mainHeight := maxY - bottomHeight - 1

	// Calculate side panel heights (divide evenly among 4 panels)
	numPanels := 4
	panelHeight := mainHeight / numPanels

	// Layout side panels
	sidePanelViews := []*gocui.View{
		gui.Views.Containers,
		gui.Views.Images,
		gui.Views.Volumes,
		gui.Views.Networks,
	}

	for i, view := range sidePanelViews {
		y0 := i * panelHeight
		y1 := (i+1)*panelHeight - 1
		if i == numPanels-1 {
			y1 = mainHeight - 1 // Last panel takes remaining space
		}

		if gui.State.ScreenMode == SCREEN_FULL {
			view.Visible = false
			continue
		}

		view.Visible = true
		_, err := g.SetView(view.Name(), 0, y0, sidePanelWidth-1, y1, 0)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
	}

	// Layout main view
	mainX0 := sidePanelWidth
	if gui.State.ScreenMode == SCREEN_FULL {
		mainX0 = 0
	}

	gui.Views.Main.Visible = true
	_, err := g.SetView("main", mainX0, 0, maxX-1, mainHeight-1, 0)
	if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
		return err
	}

	// Layout bottom bar (options)
	gui.Views.Options.Visible = true
	_, err = g.SetView("options", 0, mainHeight, maxX-1, maxY-1, 0)
	if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
		return err
	}

	// Hide unused views
	gui.Views.Information.Visible = false
	gui.Views.AppStatus.Visible = false
	gui.Views.Filter.Visible = false
	gui.Views.FilterPrefix.Visible = false

	// Render options bar
	gui.renderOptions()

	return nil
}

func (gui *Gui) showLimitView(g *gocui.Gui, maxX, maxY int) {
	gui.Views.Limit.Visible = true
	_, _ = g.SetView("limit", 0, 0, maxX-1, maxY-1, 0)
	gui.Views.Limit.Clear()
	fmt.Fprintln(gui.Views.Limit, gui.Tr.NotEnoughSpace)
	fmt.Fprintf(gui.Views.Limit, "\nMinimum size: 60x15\nCurrent size: %dx%d", maxX, maxY)

	// Hide all other views
	for _, view := range gui.getSidePanelViews() {
		view.Visible = false
	}
	gui.Views.Main.Visible = false
	gui.Views.Options.Visible = false
}

func (gui *Gui) hideLimitView() {
	gui.Views.Limit.Visible = false
}

func (gui *Gui) renderOptions() {
	if gui.Views.Options == nil {
		return
	}
	gui.Views.Options.Clear()

	// Show confirmation hints when confirmation panel is focused
	if gui.currentViewName() == "confirmation" {
		fmt.Fprint(gui.Views.Options, " y/Enter: "+gui.Tr.Yes+" │ n/Esc: "+gui.Tr.No)
		return
	}

	hints := []string{
		"↑↓/jk:navigate",
		"←→/hl:panels",
		"Enter:start",
		"s:stop",
		"d:delete",
		"[:prev tab",
		"]:next tab",
		"q:quit",
	}

	if gui.statusMessage != "" {
		fmt.Fprint(gui.Views.Options, " "+gui.statusMessage)
	} else {
		fmt.Fprint(gui.Views.Options, " "+strings.Join(hints, " │ "))
	}
}

func (gui *Gui) setStatus(msg string) {
	gui.statusMessage = msg
	gui.renderOptions()

	go func() {
		time.Sleep(3 * time.Second)
		gui.statusMessage = ""
		gui.g.Update(func(g *gocui.Gui) error {
			gui.renderOptions()
			return nil
		})
	}()
}

// IGui interface implementations

func (gui *Gui) HandleClick(v *gocui.View, itemCount int, selectedLine *int, handleSelect func() error) error {
	if gui.popupPanelFocused() && v != nil && !gui.isPopupPanel(v.Name()) {
		return nil
	}

	_, cy := v.Cursor()
	_, oy := v.Origin()

	newSelectedLine := cy + oy

	if newSelectedLine < 0 {
		newSelectedLine = 0
	}

	if newSelectedLine > itemCount-1 {
		newSelectedLine = itemCount - 1
	}

	*selectedLine = newSelectedLine

	if gui.currentViewName() != v.Name() {
		if err := gui.switchFocus(v); err != nil {
			return err
		}
	}

	return handleSelect()
}

func (gui *Gui) NewSimpleRenderStringTask(getContent func() string) tasks.TaskFunc {
	return func(ctx context.Context) {
		content := getContent()
		gui.RenderStringMain(content)
	}
}

func (gui *Gui) FocusY(selectedY int, lineCount int, v *gocui.View) {
	gui.focusPoint(0, selectedY, lineCount, v)
}

func (gui *Gui) focusPoint(selectedX int, selectedY int, lineCount int, v *gocui.View) {
	if selectedY < 0 || selectedY > lineCount {
		return
	}
	ox, oy := v.Origin()
	originalOy := oy
	cx, cy := v.Cursor()
	originalCy := cy
	_, height := v.Size()

	ly := utils.Max(height-1, 0)

	windowStart := oy
	windowEnd := oy + ly

	if selectedY < windowStart {
		oy = utils.Max(oy-(windowStart-selectedY), 0)
	} else if selectedY > windowEnd {
		oy += (selectedY - windowEnd)
	}

	if windowEnd > lineCount-1 {
		shiftAmount := (windowEnd - (lineCount - 1))
		oy = utils.Max(oy-shiftAmount, 0)
	}

	if originalOy != oy {
		_ = v.SetOrigin(ox, oy)
	}

	cy = selectedY - oy
	if originalCy != cy {
		_ = v.SetCursor(cx, selectedY-oy)
	}
}

func (gui *Gui) ShouldRefresh(key string) bool {
	if gui.State.Panels.Main.ObjectKey == key {
		return false
	}
	gui.State.Panels.Main.ObjectKey = key
	return true
}

func (gui *Gui) GetMainView() *gocui.View {
	return gui.Views.Main
}

func (gui *Gui) IsCurrentView(view *gocui.View) bool {
	return view == gui.g.CurrentView()
}

func (gui *Gui) FilterString(view *gocui.View) string {
	if gui.State.Filter.panel != nil && gui.State.Filter.panel.GetView() != view {
		return ""
	}
	return gui.State.Filter.needle
}

func (gui *Gui) IgnoreStrings() []string {
	return []string{}
}

func (gui *Gui) Update(f func() error) {
	gui.g.Update(func(*gocui.Gui) error { return f() })
}

func (gui *Gui) QueueTask(f func(ctx context.Context)) error {
	return gui.taskManager.NewTask(f)
}

func (gui *Gui) RenderStringMain(s string) {
	_ = gui.renderString(gui.g, "main", s)
}

func (gui *Gui) renderString(g *gocui.Gui, viewName, s string) error {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			return nil
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		v.Clear()
		fmt.Fprint(v, utils.NormalizeLinefeeds(s))
		return nil
	})
	return nil
}

// View helper methods

func (gui *Gui) currentViewName() string {
	currentView := gui.g.CurrentView()
	if currentView == nil {
		return gui.initiallyFocusedViewName()
	}
	return currentView.Name()
}

func (gui *Gui) isPopupPanel(viewName string) bool {
	for _, name := range gui.popupViewNames() {
		if name == viewName {
			return true
		}
	}
	return false
}

func (gui *Gui) popupPanelFocused() bool {
	return gui.isPopupPanel(gui.currentViewName())
}

func (gui *Gui) switchFocus(newView *gocui.View) error {
	if _, err := gui.g.SetCurrentView(newView.Name()); err != nil {
		return err
	}

	gui.State.ViewStack = append(gui.State.ViewStack, newView.Name())

	return gui.newLineFocused(newView)
}

func (gui *Gui) returnFocus() error {
	if len(gui.State.ViewStack) <= 1 {
		return nil
	}

	// Get the previous view (second to last in stack)
	previousViewName := gui.State.ViewStack[len(gui.State.ViewStack)-2]

	// Pop the current view from stack
	gui.State.ViewStack = gui.State.ViewStack[:len(gui.State.ViewStack)-1]

	previousView, err := gui.g.View(previousViewName)
	if err != nil {
		return err
	}

	if _, err := gui.g.SetCurrentView(previousViewName); err != nil {
		return err
	}

	return gui.newLineFocused(previousView)
}

func (gui *Gui) newLineFocused(v *gocui.View) error {
	if v == nil {
		return nil
	}

	switch v.Name() {
	case "containers":
		return gui.Panels.Containers.HandleSelect()
	case "images":
		return gui.Panels.Images.HandleSelect()
	case "volumes":
		return gui.Panels.Volumes.HandleSelect()
	case "networks":
		return gui.Panels.Networks.HandleSelect()
	case "menu":
		return gui.Panels.Menu.HandleSelect()
	}

	return nil
}

func (gui *Gui) sideViewNames() []string {
	names := []string{}
	for _, panel := range gui.allSidePanels() {
		if !panel.IsHidden() {
			names = append(names, panel.GetView().Name())
		}
	}
	return names
}

func (gui *Gui) allSidePanels() []panels.ISideListPanel {
	return []panels.ISideListPanel{
		gui.Panels.Containers,
		gui.Panels.Images,
		gui.Panels.Volumes,
		gui.Panels.Networks,
	}
}

func (gui *Gui) currentSidePanel() (panels.ISideListPanel, bool) {
	viewName := gui.currentViewName()
	for _, sidePanel := range gui.allSidePanels() {
		if sidePanel.GetView().Name() == viewName {
			return sidePanel, true
		}
	}
	return nil, false
}

func (gui *Gui) resetMainView() {
	gui.State.Panels.Main.ObjectKey = ""
	gui.Views.Main.Wrap = true
}
