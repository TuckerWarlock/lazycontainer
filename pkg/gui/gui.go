package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/config"
)

const (
	containersView       = "containers"
	mainView             = "main"
	statusView           = "status"
	UNKNOWN_VIEW_ERROR_MSG = "unknown view"
)

// Gui wraps the gocui Gui object
type Gui struct {
	g                *gocui.Gui
	Log              *logrus.Entry
	ContainerCommand *commands.ContainerCommand
	OSCommand        *commands.OSCommand
	Config           *config.AppConfig

	containers     []commands.Container
	selectedIdx    int
	showAll        bool
	mainContent    string
	statusMessage  string
	lastRefresh    time.Time
}

// NewGui creates a new Gui instance
func NewGui(log *logrus.Entry, containerCmd *commands.ContainerCommand, osCmd *commands.OSCommand, config *config.AppConfig) (*Gui, error) {
	return &Gui{
		Log:              log,
		ContainerCommand: containerCmd,
		OSCommand:        osCmd,
		Config:           config,
		showAll:          true,
		containers:       []commands.Container{},
		selectedIdx:      0,
	}, nil
}

// Run starts the GUI main loop
func (gui *Gui) Run() error {
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

	// Create all views BEFORE setting up the manager
	if err := gui.createAllViews(); err != nil {
		return fmt.Errorf("failed to create views: %w", err)
	}

	g.SetManagerFunc(gui.layout)

	if err := gui.keybindings(g); err != nil {
		return fmt.Errorf("failed to set keybindings: %w", err)
	}

	// Initial data load
	gui.loadContainers()

	// Background refresh
	go func() {
		time.Sleep(500 * time.Millisecond)
		gui.backgroundRefresh()
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

// createAllViews creates all views before the main loop
func (gui *Gui) createAllViews() error {
	// Create views with arbitrary small dimensions - they'll be resized in layout()
	viewNames := []string{containersView, mainView, statusView}

	for _, name := range viewNames {
		_, err := gui.g.SetView(name, 0, 0, 10, 10, 0)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return fmt.Errorf("failed to create view %s: %w", name, err)
		}
	}

	// Configure the views
	if v, err := gui.g.View(containersView); err == nil {
		v.Title = " Containers "
		v.Highlight = true
		v.SelBgColor = gocui.ColorBlue
		v.SelFgColor = gocui.ColorWhite
	}

	if v, err := gui.g.View(mainView); err == nil {
		v.Title = " Details "
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := gui.g.View(statusView); err == nil {
		v.Frame = true
		v.Title = " Status "
	}

	// Set the initial current view
	if _, err := gui.g.SetCurrentView(containersView); err != nil {
		return fmt.Errorf("failed to set current view: %w", err)
	}

	return nil
}

func (gui *Gui) backgroundRefresh() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		gui.refreshContainers()
		gui.g.Update(func(g *gocui.Gui) error {
			return gui.renderContainers()
		})
	}
}

// loadContainers fetches containers without updating status (for initial load)
func (gui *Gui) loadContainers() {
	containers, err := gui.ContainerCommand.ListContainers(gui.showAll)
	if err != nil {
		gui.Log.Error(err)
		return
	}
	gui.containers = containers
	gui.lastRefresh = time.Now()

	if gui.selectedIdx >= len(gui.containers) && len(gui.containers) > 0 {
		gui.selectedIdx = len(gui.containers) - 1
	}
}

func (gui *Gui) refreshContainers() {
	containers, err := gui.ContainerCommand.ListContainers(gui.showAll)
	if err != nil {
		gui.setStatus(fmt.Sprintf("Error: %v", err))
		return
	}
	gui.containers = containers
	gui.lastRefresh = time.Now()

	if gui.selectedIdx >= len(gui.containers) && len(gui.containers) > 0 {
		gui.selectedIdx = len(gui.containers) - 1
	}
}

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Minimum terminal size check
	if maxX < 50 || maxY < 10 {
		// Hide normal views
		if v, err := g.View(containersView); err == nil {
			v.Visible = false
		}
		if v, err := g.View(mainView); err == nil {
			v.Visible = false
		}
		if v, err := g.View(statusView); err == nil {
			v.Visible = false
		}
		return nil
	}

	// Calculate dimensions
	listWidth := maxX * 40 / 100
	if listWidth < 25 {
		listWidth = 25
	}
	if listWidth > maxX-25 {
		listWidth = maxX - 25
	}

	mainHeight := maxY - 4
	if mainHeight < 5 {
		mainHeight = 5
	}

	// Resize container list view (left panel)
	_, err := g.SetView(containersView, 0, 0, listWidth-1, mainHeight, 0)
	if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
		return fmt.Errorf("containers view SetView: %w", err)
	}
	if v, err := g.View(containersView); err == nil {
		v.Visible = true
	}

	// Resize main content view (right panel)
	_, err = g.SetView(mainView, listWidth, 0, maxX-1, mainHeight, 0)
	if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
		return fmt.Errorf("main view SetView: %w", err)
	}
	if v, err := g.View(mainView); err == nil {
		v.Visible = true
	}

	// Resize status bar (bottom)
	_, err = g.SetView(statusView, 0, mainHeight+1, maxX-1, maxY-1, 0)
	if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
		return fmt.Errorf("status view SetView: %w", err)
	}
	if v, err := g.View(statusView); err == nil {
		v.Visible = true
	}

	// Render content
	if err := gui.renderContainers(); err != nil {
		return fmt.Errorf("renderContainers: %w", err)
	}
	if err := gui.renderMain(); err != nil {
		return fmt.Errorf("renderMain: %w", err)
	}
	if err := gui.renderStatus(); err != nil {
		return fmt.Errorf("renderStatus: %w", err)
	}

	return nil
}

func (gui *Gui) renderContainers() error {
	if gui.g == nil {
		return nil
	}
	v, err := gui.g.View(containersView)
	if err != nil {
		return nil // View might not exist yet
	}
	v.Clear()

	if len(gui.containers) == 0 {
		fmt.Fprintln(v, "  No containers found")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "  Press 'a' to toggle showing")
		fmt.Fprintln(v, "  all containers (including stopped)")
		return nil
	}

	for i, container := range gui.containers {
		line := gui.formatContainerLine(container, i == gui.selectedIdx)
		fmt.Fprintln(v, line)
	}

	return nil
}

func (gui *Gui) formatContainerLine(c commands.Container, selected bool) string {
	// Status indicator with color
	var statusIndicator string
	if c.IsRunning() {
		statusIndicator = color.GreenString("●")
	} else {
		statusIndicator = color.RedString("●")
	}

	name := c.GetName()
	if len(name) > 20 {
		name = name[:17] + "..."
	}

	status := c.GetStatus()
	uptime := c.GetUptime()

	// Format: ● name          status  uptime
	line := fmt.Sprintf(" %s %-20s %-8s %s", statusIndicator, name, status, uptime)

	return line
}

func (gui *Gui) renderMain() error {
	if gui.g == nil {
		return nil
	}
	v, err := gui.g.View(mainView)
	if err != nil {
		return nil // View might not exist yet
	}
	v.Clear()

	if len(gui.containers) == 0 {
		fmt.Fprintln(v, "No container selected")
		return nil
	}

	if gui.selectedIdx >= len(gui.containers) {
		return nil
	}

	container := gui.containers[gui.selectedIdx]

	// Container details
	fmt.Fprintf(v, "Name:     %s\n", color.CyanString(container.GetName()))
	fmt.Fprintf(v, "Status:   %s\n", gui.colorStatus(container.GetStatus()))
	fmt.Fprintf(v, "Image:    %s\n", container.GetImage())
	fmt.Fprintf(v, "CPUs:     %d\n", container.GetCPUs())
	fmt.Fprintf(v, "Memory:   %s\n", container.GetMemoryHuman())

	if ports := container.GetPorts(); ports != "" {
		fmt.Fprintf(v, "Ports:    %s\n", ports)
	}

	if container.IsRunning() {
		fmt.Fprintf(v, "Uptime:   %s\n", container.GetUptime())
		fmt.Fprintf(v, "Started:  %s\n", container.GetStartedAt().Format("2006-01-02 15:04:05"))
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "─────────────────────────────────────")
	fmt.Fprintln(v, "")

	// Show recent logs
	if container.IsRunning() {
		logs, err := gui.ContainerCommand.GetContainerLogs(container.GetID(), 20)
		if err == nil && logs != "" {
			fmt.Fprintln(v, color.YellowString("Recent Logs:"))
			fmt.Fprintln(v, logs)
		}
	} else {
		fmt.Fprintln(v, color.YellowString("Container is not running"))
		fmt.Fprintln(v, "Press Enter to start")
	}

	return nil
}

func (gui *Gui) colorStatus(status string) string {
	switch status {
	case "running":
		return color.GreenString(status)
	case "stopped", "exited":
		return color.RedString(status)
	case "paused":
		return color.YellowString(status)
	default:
		return status
	}
}

func (gui *Gui) renderStatus() error {
	if gui.g == nil {
		return nil
	}
	v, err := gui.g.View(statusView)
	if err != nil {
		return nil // View might not exist yet
	}
	v.Clear()

	// Key hints
	hints := []string{
		"↑↓:navigate",
		"Enter:start",
		"s:stop",
		"d:delete",
		"l:logs",
		"a:all",
		"PgUp/Dn:scroll",
		"q:quit",
	}

	statusLine := " " + strings.Join(hints, " │ ")

	if gui.statusMessage != "" {
		statusLine = " " + gui.statusMessage
	}

	fmt.Fprint(v, statusLine)

	return nil
}

func (gui *Gui) setStatus(msg string) {
	gui.statusMessage = msg
	go func() {
		time.Sleep(3 * time.Second)
		gui.statusMessage = ""
		gui.g.Update(func(g *gocui.Gui) error {
			return gui.renderStatus()
		})
	}()
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	// Global keybindings
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		return err
	}

	// Scroll main view with page up/down (global)
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, gui.scrollMainUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, gui.scrollMainDown); err != nil {
		return err
	}

	// Container view specific keybindings
	bindings := []struct {
		key     interface{}
		handler func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyArrowUp, gui.cursorUp},
		{gocui.KeyArrowDown, gui.cursorDown},
		{'k', gui.cursorUp},
		{'j', gui.cursorDown},
		{gocui.KeyEnter, gui.startContainer},
		{'s', gui.stopContainer},
		{'d', gui.deleteContainer},
		{'l', gui.showLogs},
		{'a', gui.toggleShowAll},
		{'r', gui.forceRefresh},
		{'q', gui.quit},
	}

	for _, b := range bindings {
		if err := g.SetKeybinding(containersView, b.key, gocui.ModNone, b.handler); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (gui *Gui) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if gui.selectedIdx > 0 {
		gui.selectedIdx--
	}
	return gui.updateView()
}

func (gui *Gui) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if gui.selectedIdx < len(gui.containers)-1 {
		gui.selectedIdx++
	}
	return gui.updateView()
}

func (gui *Gui) updateView() error {
	if err := gui.renderContainers(); err != nil {
		return err
	}
	return gui.renderMain()
}

func (gui *Gui) getSelectedContainer() *commands.Container {
	if gui.selectedIdx >= 0 && gui.selectedIdx < len(gui.containers) {
		return &gui.containers[gui.selectedIdx]
	}
	return nil
}

func (gui *Gui) startContainer(g *gocui.Gui, v *gocui.View) error {
	container := gui.getSelectedContainer()
	if container == nil {
		return nil
	}

	if container.IsRunning() {
		gui.setStatus("Container is already running")
		return nil
	}

	gui.setStatus(fmt.Sprintf("Starting %s...", container.GetName()))

	go func() {
		err := gui.ContainerCommand.StartContainer(container.GetID())
		gui.g.Update(func(g *gocui.Gui) error {
			if err != nil {
				gui.setStatus(fmt.Sprintf("Error: %v", err))
			} else {
				gui.setStatus(fmt.Sprintf("Started %s", container.GetName()))
				gui.refreshContainers()
			}
			return gui.updateView()
		})
	}()

	return nil
}

func (gui *Gui) stopContainer(g *gocui.Gui, v *gocui.View) error {
	container := gui.getSelectedContainer()
	if container == nil {
		return nil
	}

	if !container.IsRunning() {
		gui.setStatus("Container is not running")
		return nil
	}

	gui.setStatus(fmt.Sprintf("Stopping %s...", container.GetName()))

	go func() {
		err := gui.ContainerCommand.StopContainer(container.GetID())
		gui.g.Update(func(g *gocui.Gui) error {
			if err != nil {
				gui.setStatus(fmt.Sprintf("Error: %v", err))
			} else {
				gui.setStatus(fmt.Sprintf("Stopped %s", container.GetName()))
				gui.refreshContainers()
			}
			return gui.updateView()
		})
	}()

	return nil
}

func (gui *Gui) deleteContainer(g *gocui.Gui, v *gocui.View) error {
	container := gui.getSelectedContainer()
	if container == nil {
		return nil
	}

	if container.IsRunning() {
		gui.setStatus("Stop container before deleting")
		return nil
	}

	gui.setStatus(fmt.Sprintf("Deleting %s...", container.GetName()))

	go func() {
		err := gui.ContainerCommand.DeleteContainer(container.GetID())
		gui.g.Update(func(g *gocui.Gui) error {
			if err != nil {
				gui.setStatus(fmt.Sprintf("Error: %v", err))
			} else {
				gui.setStatus(fmt.Sprintf("Deleted %s", container.GetName()))
				gui.refreshContainers()
			}
			return gui.updateView()
		})
	}()

	return nil
}

func (gui *Gui) showLogs(g *gocui.Gui, v *gocui.View) error {
	container := gui.getSelectedContainer()
	if container == nil {
		return nil
	}

	mainV, err := gui.g.View(mainView)
	if err != nil {
		return err
	}
	mainV.Clear()

	mainV.Title = fmt.Sprintf(" Logs: %s ", container.GetName())

	logs, err := gui.ContainerCommand.GetContainerLogs(container.GetID(), 100)
	if err != nil {
		fmt.Fprintf(mainV, "Error fetching logs: %v", err)
		return nil
	}

	fmt.Fprint(mainV, logs)
	return nil
}

func (gui *Gui) toggleShowAll(g *gocui.Gui, v *gocui.View) error {
	gui.showAll = !gui.showAll
	if gui.showAll {
		gui.setStatus("Showing all containers")
	} else {
		gui.setStatus("Showing running containers only")
	}
	gui.refreshContainers()
	return gui.updateView()
}

func (gui *Gui) forceRefresh(g *gocui.Gui, v *gocui.View) error {
	gui.setStatus("Refreshing...")
	gui.refreshContainers()
	return gui.updateView()
}

func (gui *Gui) scrollMainUp(g *gocui.Gui, v *gocui.View) error {
	mainV, err := gui.g.View(mainView)
	if err != nil {
		return nil
	}
	mainV.Autoscroll = false
	ox, oy := mainV.Origin()
	if oy > 0 {
		return mainV.SetOrigin(ox, oy-1)
	}
	return nil
}

func (gui *Gui) scrollMainDown(g *gocui.Gui, v *gocui.View) error {
	mainV, err := gui.g.View(mainView)
	if err != nil {
		return nil
	}
	mainV.Autoscroll = false
	ox, oy := mainV.Origin()
	return mainV.SetOrigin(ox, oy+1)
}
