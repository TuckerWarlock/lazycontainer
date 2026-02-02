package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/gui/panels"
)

func (gui *Gui) keybindings(g *gocui.Gui) error {
	// Global keybindings (no view name = applies everywhere)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, gui.quit); err != nil {
		return err
	}

	// Global scroll for main panel
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, gui.scrollMainUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, gui.scrollMainDown); err != nil {
		return err
	}

	// Tab navigation between panels (global)
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, gui.nextView); err != nil {
		return err
	}

	// Main tab navigation within panels (global)
	if err := g.SetKeybinding("", '[', gocui.ModNone, gui.handlePrevMainTab); err != nil {
		return err
	}
	if err := g.SetKeybinding("", ']', gocui.ModNone, gui.handleNextMainTab); err != nil {
		return err
	}

	// Number keys to focus specific panels (global - no view name)
	if err := g.SetKeybinding("", '1', gocui.ModNone, gui.focusContainersPanel); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '2', gocui.ModNone, gui.focusImagesPanel); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '3', gocui.ModNone, gui.focusVolumesPanel); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '4', gocui.ModNone, gui.focusNetworksPanel); err != nil {
		return err
	}

	// Bind navigation, mouse, and actions to all side panels
	for _, panel := range gui.allSidePanels() {
		viewName := panel.GetView().Name()
		if err := gui.bindPanelNavigationKeys(g, viewName, panel); err != nil {
			return err
		}
	}

	// Main panel bindings (scroll and click)
	if err := gui.bindMainPanelKeys(g); err != nil {
		return err
	}

	// Container-specific keybindings
	if err := g.SetKeybinding("containers", gocui.KeyEnter, gocui.ModNone, gui.handleContainerStart); err != nil {
		return err
	}
	if err := g.SetKeybinding("containers", 's', gocui.ModNone, gui.handleContainerStop); err != nil {
		return err
	}
	if err := g.SetKeybinding("containers", 'd', gocui.ModNone, gui.handleContainerDelete); err != nil {
		return err
	}
	if err := g.SetKeybinding("containers", 'a', gocui.ModNone, gui.handleToggleShowAll); err != nil {
		return err
	}
	if err := g.SetKeybinding("containers", 'e', gocui.ModNone, gui.handleToggleShowAll); err != nil {
		return err
	}

	// Image-specific keybindings
	if err := g.SetKeybinding("images", 'd', gocui.ModNone, gui.handleImageDelete); err != nil {
		return err
	}

	// Volume-specific keybindings
	if err := g.SetKeybinding("volumes", 'd', gocui.ModNone, gui.handleVolumeDelete); err != nil {
		return err
	}

	// Network-specific keybindings
	if err := g.SetKeybinding("networks", 'd', gocui.ModNone, gui.handleNetworkDelete); err != nil {
		return err
	}

	return nil
}

// bindPanelNavigationKeys sets up keyboard and mouse bindings for a side panel
func (gui *Gui) bindPanelNavigationKeys(g *gocui.Gui, viewName string, panel panels.ISideListPanel) error {
	// Keyboard navigation - up
	if err := g.SetKeybinding(viewName, gocui.KeyArrowUp, gocui.ModNone, gui.wrappedPanelHandler(panel.HandlePrevLine)); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewName, 'k', gocui.ModNone, gui.wrappedPanelHandler(panel.HandlePrevLine)); err != nil {
		return err
	}

	// Keyboard navigation - down
	if err := g.SetKeybinding(viewName, gocui.KeyArrowDown, gocui.ModNone, gui.wrappedPanelHandler(panel.HandleNextLine)); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewName, 'j', gocui.ModNone, gui.wrappedPanelHandler(panel.HandleNextLine)); err != nil {
		return err
	}

	// Panel navigation - left/right
	if err := g.SetKeybinding(viewName, gocui.KeyArrowLeft, gocui.ModNone, gui.previousView); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewName, 'h', gocui.ModNone, gui.previousView); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewName, gocui.KeyArrowRight, gocui.ModNone, gui.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewName, 'l', gocui.ModNone, gui.nextView); err != nil {
		return err
	}

	// Mouse wheel - scroll up/down
	if err := g.SetKeybinding(viewName, gocui.MouseWheelUp, gocui.ModNone, gui.wrappedPanelHandler(panel.HandlePrevLine)); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewName, gocui.MouseWheelDown, gocui.ModNone, gui.wrappedPanelHandler(panel.HandleNextLine)); err != nil {
		return err
	}

	// Mouse click - select item
	if err := g.SetKeybinding(viewName, gocui.MouseLeft, gocui.ModNone, gui.wrappedPanelHandler(panel.HandleClick)); err != nil {
		return err
	}

	// Refresh
	if err := g.SetKeybinding(viewName, 'r', gocui.ModNone, gui.handleRefresh); err != nil {
		return err
	}

	return nil
}

// bindMainPanelKeys sets up bindings for the main panel
func (gui *Gui) bindMainPanelKeys(g *gocui.Gui) error {
	// Mouse wheel scroll in main panel
	if err := g.SetKeybinding("main", gocui.MouseWheelUp, gocui.ModNone, gui.scrollMainUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.MouseWheelDown, gocui.ModNone, gui.scrollMainDown); err != nil {
		return err
	}

	// Arrow key scroll in main panel
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, gui.scrollMainUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, gui.scrollMainDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'k', gocui.ModNone, gui.scrollMainUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'j', gocui.ModNone, gui.scrollMainDown); err != nil {
		return err
	}

	// Click in main panel - just focus it
	if err := g.SetKeybinding("main", gocui.MouseLeft, gocui.ModNone, gui.handleMainClick); err != nil {
		return err
	}

	// Left/right in main panel - go back to side panels
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, gui.previousView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'h', gocui.ModNone, gui.previousView); err != nil {
		return err
	}

	return nil
}

// wrappedPanelHandler wraps a panel handler to be used as a gocui handler
func (gui *Gui) wrappedPanelHandler(handler func() error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return handler()
	}
}

// Global handlers

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (gui *Gui) scrollMainUp(g *gocui.Gui, v *gocui.View) error {
	mainV := gui.Views.Main
	if mainV == nil {
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
	mainV := gui.Views.Main
	if mainV == nil {
		return nil
	}
	mainV.Autoscroll = false
	ox, oy := mainV.Origin()
	return mainV.SetOrigin(ox, oy+1)
}

func (gui *Gui) handleMainClick(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}
	return nil
}

// Panel navigation

func (gui *Gui) nextView(g *gocui.Gui, v *gocui.View) error {
	sideViewNames := gui.sideViewNames()
	if len(sideViewNames) == 0 {
		return nil
	}

	var focusedViewName string
	currentName := ""
	if v != nil {
		currentName = v.Name()
	}

	// If we're in main panel or at the end, go to first side panel
	if currentName == "main" || currentName == sideViewNames[len(sideViewNames)-1] {
		focusedViewName = sideViewNames[0]
	} else {
		// Find current position and go to next
		for i, name := range sideViewNames {
			if currentName == name {
				focusedViewName = sideViewNames[i+1]
				break
			}
		}
	}

	if focusedViewName == "" {
		focusedViewName = sideViewNames[0]
	}

	focusedView, err := g.View(focusedViewName)
	if err != nil {
		return err
	}
	gui.resetMainView()
	return gui.switchFocus(focusedView)
}

func (gui *Gui) previousView(g *gocui.Gui, v *gocui.View) error {
	sideViewNames := gui.sideViewNames()
	if len(sideViewNames) == 0 {
		return nil
	}

	var focusedViewName string
	currentName := ""
	if v != nil {
		currentName = v.Name()
	}

	// If we're in main panel or at the start, go to last side panel
	if currentName == "main" || currentName == sideViewNames[0] {
		focusedViewName = sideViewNames[len(sideViewNames)-1]
	} else {
		// Find current position and go to previous
		for i, name := range sideViewNames {
			if currentName == name {
				focusedViewName = sideViewNames[i-1]
				break
			}
		}
	}

	if focusedViewName == "" {
		focusedViewName = sideViewNames[len(sideViewNames)-1]
	}

	focusedView, err := g.View(focusedViewName)
	if err != nil {
		return err
	}
	gui.resetMainView()
	return gui.switchFocus(focusedView)
}

func (gui *Gui) focusContainersPanel(g *gocui.Gui, v *gocui.View) error {
	return gui.focusPanel("containers")
}

func (gui *Gui) focusImagesPanel(g *gocui.Gui, v *gocui.View) error {
	return gui.focusPanel("images")
}

func (gui *Gui) focusVolumesPanel(g *gocui.Gui, v *gocui.View) error {
	return gui.focusPanel("volumes")
}

func (gui *Gui) focusNetworksPanel(g *gocui.Gui, v *gocui.View) error {
	return gui.focusPanel("networks")
}

func (gui *Gui) focusPanel(viewName string) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return err
	}
	gui.resetMainView()
	return gui.switchFocus(view)
}

// Main tab navigation

func (gui *Gui) handlePrevMainTab(g *gocui.Gui, v *gocui.View) error {
	panel, ok := gui.currentSidePanel()
	if !ok {
		return nil
	}
	return panel.HandlePrevMainTab()
}

func (gui *Gui) handleNextMainTab(g *gocui.Gui, v *gocui.View) error {
	panel, ok := gui.currentSidePanel()
	if !ok {
		return nil
	}
	return panel.HandleNextMainTab()
}

// Line navigation (kept for compatibility but now handled via wrapped handlers)

func (gui *Gui) handlePrevLine(g *gocui.Gui, v *gocui.View) error {
	panel, ok := gui.currentSidePanel()
	if !ok {
		return nil
	}
	return panel.HandlePrevLine()
}

func (gui *Gui) handleNextLine(g *gocui.Gui, v *gocui.View) error {
	panel, ok := gui.currentSidePanel()
	if !ok {
		return nil
	}
	return panel.HandleNextLine()
}

// Refresh

func (gui *Gui) handleRefresh(g *gocui.Gui, v *gocui.View) error {
	gui.setStatus("Refreshing...")
	gui.refresh()
	return nil
}

// Container actions

func (gui *Gui) handleContainerStart(g *gocui.Gui, v *gocui.View) error {
	container, err := gui.Panels.Containers.GetSelectedItem()
	if err != nil {
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
				gui.refresh()
			}
			return nil
		})
	}()

	return nil
}

func (gui *Gui) handleContainerStop(g *gocui.Gui, v *gocui.View) error {
	container, err := gui.Panels.Containers.GetSelectedItem()
	if err != nil {
		return nil
	}

	if !container.IsRunning() {
		gui.setStatus("Container is not running")
		return nil
	}

	return gui.createConfirmationPanel(
		gui.Tr.Confirm,
		gui.Tr.StopContainer,
		func(g *gocui.Gui, v *gocui.View) error {
			gui.setStatus(fmt.Sprintf("Stopping %s...", container.GetName()))

			go func() {
				err := gui.ContainerCommand.StopContainer(container.GetID())
				gui.g.Update(func(g *gocui.Gui) error {
					if err != nil {
						gui.setStatus(fmt.Sprintf("Error: %v", err))
					} else {
						gui.setStatus(fmt.Sprintf("Stopped %s", container.GetName()))
						gui.refresh()
					}
					return nil
				})
			}()

			return nil
		},
		nil,
	)
}

func (gui *Gui) handleContainerDelete(g *gocui.Gui, v *gocui.View) error {
	container, err := gui.Panels.Containers.GetSelectedItem()
	if err != nil {
		return nil
	}

	if container.IsRunning() {
		gui.setStatus("Stop container before deleting")
		return nil
	}

	return gui.createConfirmationPanel(
		gui.Tr.Confirm,
		gui.Tr.ConfirmRemoveContainer,
		func(g *gocui.Gui, v *gocui.View) error {
			gui.setStatus(fmt.Sprintf("Deleting %s...", container.GetName()))

			go func() {
				err := gui.ContainerCommand.DeleteContainer(container.GetID())
				gui.g.Update(func(g *gocui.Gui) error {
					if err != nil {
						gui.setStatus(fmt.Sprintf("Error: %v", err))
					} else {
						gui.setStatus(fmt.Sprintf("Deleted %s", container.GetName()))
						gui.refresh()
					}
					return nil
				})
			}()

			return nil
		},
		nil,
	)
}

func (gui *Gui) handleToggleShowAll(g *gocui.Gui, v *gocui.View) error {
	gui.State.ShowExitedContainers = !gui.State.ShowExitedContainers
	if gui.State.ShowExitedContainers {
		gui.setStatus("Showing all containers")
	} else {
		gui.setStatus("Showing running containers only")
	}
	return gui.refreshContainers()
}

// Image actions

func (gui *Gui) handleImageDelete(g *gocui.Gui, v *gocui.View) error {
	image, err := gui.Panels.Images.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.createConfirmationPanel(
		gui.Tr.Confirm,
		gui.Tr.ConfirmRemoveImage,
		func(g *gocui.Gui, v *gocui.View) error {
			gui.setStatus(fmt.Sprintf("Deleting image %s...", image.Reference))

			go func() {
				err := gui.ContainerCommand.DeleteImage(image.Reference)
				gui.g.Update(func(g *gocui.Gui) error {
					if err != nil {
						gui.setStatus(fmt.Sprintf("Error: %v", err))
					} else {
						gui.setStatus(fmt.Sprintf("Deleted image %s", image.Reference))
						gui.refresh()
					}
					return nil
				})
			}()

			return nil
		},
		nil,
	)
}

// Volume actions

func (gui *Gui) handleVolumeDelete(g *gocui.Gui, v *gocui.View) error {
	volume, err := gui.Panels.Volumes.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.createConfirmationPanel(
		gui.Tr.Confirm,
		gui.Tr.ConfirmRemoveVolume,
		func(g *gocui.Gui, v *gocui.View) error {
			gui.setStatus(fmt.Sprintf("Deleting volume %s...", volume.Name))

			go func() {
				err := gui.ContainerCommand.DeleteVolume(volume.Name)
				gui.g.Update(func(g *gocui.Gui) error {
					if err != nil {
						gui.setStatus(fmt.Sprintf("Error: %v", err))
					} else {
						gui.setStatus(fmt.Sprintf("Deleted volume %s", volume.Name))
						gui.refresh()
					}
					return nil
				})
			}()

			return nil
		},
		nil,
	)
}

// Network actions

func (gui *Gui) handleNetworkDelete(g *gocui.Gui, v *gocui.View) error {
	network, err := gui.Panels.Networks.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.createConfirmationPanel(
		gui.Tr.Confirm,
		gui.Tr.ConfirmRemoveNetwork,
		func(g *gocui.Gui, v *gocui.View) error {
			gui.setStatus(fmt.Sprintf("Deleting network %s...", network.ID))

			go func() {
				err := gui.ContainerCommand.DeleteNetwork(network.ID)
				gui.g.Update(func(g *gocui.Gui) error {
					if err != nil {
						gui.setStatus(fmt.Sprintf("Error: %v", err))
					} else {
						gui.setStatus(fmt.Sprintf("Deleted network %s", network.ID))
						gui.refresh()
					}
					return nil
				})
			}()

			return nil
		},
		nil,
	)
}

// Helper to get selected container (for backward compatibility)
func (gui *Gui) getSelectedContainer() *commands.Container {
	container, err := gui.Panels.Containers.GetSelectedItem()
	if err != nil {
		return nil
	}
	return container
}

// Helper to get current list panel
func (gui *Gui) currentListPanel() (panels.ISideListPanel, bool) {
	return gui.currentSidePanel()
}
