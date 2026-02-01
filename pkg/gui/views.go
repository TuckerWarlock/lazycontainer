package gui

import (
	"github.com/jesseduffield/gocui"
)

// Views holds all the gocui views for the GUI
type Views struct {
	// side panels
	Containers *gocui.View
	Images     *gocui.View
	Volumes    *gocui.View
	Networks   *gocui.View

	// main panel
	Main *gocui.View

	// bottom line
	Options      *gocui.View
	Information  *gocui.View
	AppStatus    *gocui.View
	FilterPrefix *gocui.View
	Filter       *gocui.View

	// popups
	Confirmation *gocui.View
	Menu         *gocui.View

	// error view for insufficient space
	Limit *gocui.View
}

type viewNameMapping struct {
	viewPtr      **gocui.View
	name         string
	autoPosition bool
}

func (gui *Gui) orderedViewNameMappings() []viewNameMapping {
	return []viewNameMapping{
		// Side panels (order matters for tab navigation)
		{viewPtr: &gui.Views.Containers, name: "containers", autoPosition: true},
		{viewPtr: &gui.Views.Images, name: "images", autoPosition: true},
		{viewPtr: &gui.Views.Volumes, name: "volumes", autoPosition: true},
		{viewPtr: &gui.Views.Networks, name: "networks", autoPosition: true},

		// Main panel
		{viewPtr: &gui.Views.Main, name: "main", autoPosition: true},

		// Bottom line
		{viewPtr: &gui.Views.Options, name: "options", autoPosition: true},
		{viewPtr: &gui.Views.AppStatus, name: "appStatus", autoPosition: true},
		{viewPtr: &gui.Views.Information, name: "information", autoPosition: true},
		{viewPtr: &gui.Views.Filter, name: "filter", autoPosition: true},
		{viewPtr: &gui.Views.FilterPrefix, name: "filterPrefix", autoPosition: true},

		// Popups (order matters - later items render on top)
		{viewPtr: &gui.Views.Menu, name: "menu", autoPosition: false},
		{viewPtr: &gui.Views.Confirmation, name: "confirmation", autoPosition: false},

		// Full-screen error view
		{viewPtr: &gui.Views.Limit, name: "limit", autoPosition: true},
	}
}

func (gui *Gui) prepareView(name string) (*gocui.View, error) {
	// Create view with minimal dimensions - will be resized in layout
	view, err := gui.g.SetView(name, 0, 0, 10, 10, 0)
	if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
		return nil, err
	}
	return view, nil
}

func (gui *Gui) initializeViews() error {
	gui.Views = Views{}

	frameRunes := []rune{'─', '│', '╭', '╮', '╰', '╯'}

	var err error
	for _, mapping := range gui.orderedViewNameMappings() {
		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		(*mapping.viewPtr).FrameRunes = frameRunes
		(*mapping.viewPtr).FgColor = gocui.ColorDefault
	}

	selectedLineBgColor := gocui.ColorBlue

	// Configure Main view
	gui.Views.Main.Wrap = true
	gui.Views.Main.IgnoreCarriageReturns = true
	gui.Views.Main.Title = " Details "

	// Configure side panels
	gui.Views.Containers.Highlight = true
	gui.Views.Containers.Title = gui.Tr.ContainersTitle
	gui.Views.Containers.TitlePrefix = "[1]"
	gui.Views.Containers.SelBgColor = selectedLineBgColor
	gui.Views.Containers.SelFgColor = gocui.ColorWhite

	gui.Views.Images.Highlight = true
	gui.Views.Images.Title = gui.Tr.ImagesTitle
	gui.Views.Images.TitlePrefix = "[2]"
	gui.Views.Images.SelBgColor = selectedLineBgColor
	gui.Views.Images.SelFgColor = gocui.ColorWhite

	gui.Views.Volumes.Highlight = true
	gui.Views.Volumes.Title = gui.Tr.VolumesTitle
	gui.Views.Volumes.TitlePrefix = "[3]"
	gui.Views.Volumes.SelBgColor = selectedLineBgColor
	gui.Views.Volumes.SelFgColor = gocui.ColorWhite

	gui.Views.Networks.Highlight = true
	gui.Views.Networks.Title = gui.Tr.NetworksTitle
	gui.Views.Networks.TitlePrefix = "[4]"
	gui.Views.Networks.SelBgColor = selectedLineBgColor
	gui.Views.Networks.SelFgColor = gocui.ColorWhite

	// Configure bottom line views
	gui.Views.Options.Frame = false
	gui.Views.Options.FgColor = gocui.ColorCyan

	gui.Views.AppStatus.FgColor = gocui.ColorCyan
	gui.Views.AppStatus.Frame = false

	gui.Views.Information.Frame = false
	gui.Views.Information.FgColor = gocui.ColorGreen

	// Configure popups
	gui.Views.Confirmation.Visible = false
	gui.Views.Confirmation.Wrap = true

	gui.Views.Menu.Visible = false
	gui.Views.Menu.SelBgColor = selectedLineBgColor

	// Configure error view
	gui.Views.Limit.Visible = false
	gui.Views.Limit.Title = gui.Tr.NotEnoughSpace
	gui.Views.Limit.Wrap = true

	// Configure filter views
	gui.Views.FilterPrefix.BgColor = gocui.ColorDefault
	gui.Views.FilterPrefix.FgColor = gocui.ColorGreen
	gui.Views.FilterPrefix.Frame = false

	gui.Views.Filter.BgColor = gocui.ColorDefault
	gui.Views.Filter.FgColor = gocui.ColorGreen
	gui.Views.Filter.Editable = true
	gui.Views.Filter.Frame = false

	return nil
}

// getSidePanelViews returns the views that appear in the side panel
func (gui *Gui) getSidePanelViews() []*gocui.View {
	return []*gocui.View{
		gui.Views.Containers,
		gui.Views.Images,
		gui.Views.Volumes,
		gui.Views.Networks,
	}
}

// getSidePanelNames returns the names of the side panel views
func (gui *Gui) getSidePanelNames() []string {
	return []string{
		"containers",
		"images",
		"volumes",
		"networks",
	}
}

// popupViewNames returns the names of popup views
func (gui *Gui) popupViewNames() []string {
	return []string{"confirmation", "menu"}
}
