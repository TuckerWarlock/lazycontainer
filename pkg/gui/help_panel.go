package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
)

// getHelpContent returns help text based on the current focused panel
func (gui *Gui) getHelpContent() string {
	currentView := gui.currentViewName()

	helpText := `
  GLOBAL KEYBINDINGS
  ==================
  q, Ctrl+C      Quit
  1,2,3,4        Switch to Containers/Images/Volumes/Networks
  Tab            Cycle panels
  [, ]           Previous/next tab (logs/config/stats)
  t              Cycle themes (default/light/dark)
  r              Refresh
  
  NAVIGATION
  ==========
  j/↑, k/↓       Move up/down
  h/←, l/→       Previous/next panel
  Enter          Select item
  
  FILTERING
  =========
  /              Start filtering
  Enter          Apply filter
  Esc            Cancel filter
  
  COPY
  ====
  c              Copy selected resource ID
  C              Copy logs from main panel
`

	// Add context-specific help
	switch currentView {
	case "containers":
		helpText += `
  CONTAINERS
  ==========
  Enter          Start container
  s              Stop container (with confirmation)
  d              Delete container (with confirmation)
  e              Exec interactive shell
  f              Follow logs (streaming)
  v              Stream live stats
  Space          Toggle multi-select
  a              Toggle show all containers (running + stopped)
`
	case "images":
		helpText += `
  IMAGES
  ======
  d              Delete image (with confirmation)
`
	case "volumes":
		helpText += `
  VOLUMES
  =======
  d              Delete volume (with confirmation)
`
	case "networks":
		helpText += `
  NETWORKS
  ========
  d              Delete network (with confirmation)
`
	}

	helpText += `
  Press Esc or q to close this help
`
	return strings.TrimSpace(helpText)
}

func (gui *Gui) closeHelpPrompt() error {
	if err := gui.returnFocus(); err != nil {
		return err
	}
	gui.g.DeleteViewKeybindings("help")
	gui.Views.Help.Visible = false
	return nil
}

func (gui *Gui) getHelpPanelDimensions(helpText string) (int, int, int, int) {
	width, height := gui.g.Size()
	// Help panel takes 80% of screen width and 70% of height
	panelWidth := width * 4 / 5
	panelHeight := gui.getMessageHeight(true, helpText, panelWidth)
	if panelHeight > height*7/10 {
		panelHeight = height * 7 / 10
	}
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) prepareHelpPanel(helpText string) error {
	x0, y0, x1, y1 := gui.getHelpPanelDimensions(helpText)
	helpView := gui.Views.Help
	_, err := gui.g.SetView("help", x0, y0, x1, y1, 0)
	if err != nil {
		return err
	}
	helpView.Title = " Help "
	helpView.Visible = true
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.switchFocus(helpView)
	})
	return nil
}

func (gui *Gui) createHelpPanel() error {
	helpText := gui.getHelpContent()
	gui.onNewPopupPanel()
	gui.g.Update(func(g *gocui.Gui) error {
		if gui.currentViewName() == "help" {
			if err := gui.closeHelpPrompt(); err != nil {
				gui.Log.Error(err.Error())
			}
		}
		err := gui.prepareHelpPanel(helpText)
		if err != nil {
			return err
		}
		gui.Views.Help.Editable = false
		if err := gui.renderString(g, "help", helpText); err != nil {
			return err
		}
		return gui.setHelpKeyBindings(g)
	})
	return nil
}

func (gui *Gui) setHelpKeyBindings(g *gocui.Gui) error {
	// Close with Esc
	if err := g.SetKeybinding("help", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gui.closeHelpPrompt()
	}); err != nil {
		return err
	}

	// Close with q
	if err := g.SetKeybinding("help", 'q', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gui.closeHelpPrompt()
	}); err != nil {
		return err
	}

	// Close with Enter
	if err := g.SetKeybinding("help", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gui.closeHelpPrompt()
	}); err != nil {
		return err
	}

	// Mouse wheel scroll
	if err := g.SetKeybinding("help", gocui.MouseWheelUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_ = v.SetCursor(0, 0)
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.MouseWheelDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		// Allow scrolling down
		ox, oy := v.Origin()
		_ = v.SetOrigin(ox, oy+1)
		return nil
	}); err != nil {
		return err
	}

	return nil
}
