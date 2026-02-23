package panels

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/warl0ck/lazycontainer/pkg/tasks"
	"github.com/warl0ck/lazycontainer/pkg/utils"
	"github.com/samber/lo"
)

type ISideListPanel interface {
	SetMainTabIndex(int)
	HandleSelect() error
	GetView() *gocui.View
	Refocus()
	RerenderList() error
	IsFilterDisabled() bool
	IsHidden() bool
	HandleNextLine() error
	HandlePrevLine() error
	HandleClick() error
	HandlePrevMainTab() error
	HandleNextMainTab() error
}

// list panel at the side of the screen that renders content to the main panel
type SideListPanel[T comparable] struct {
	ContextState *ContextState[T]

	ListPanel[T]

	// message to render in the main view if there are no items in the panel
	// and it has focus. Leave empty if you don't want to render anything
	NoItemsMessage string

	// a representation of the gui
	Gui IGui

	// this Filter is applied on top of additional default filters
	Filter func(T) bool
	Sort   func(a, b T) bool

	// a callback to invoke when the item is clicked
	OnClick func(T) error

	// returns the cells that we render to the view in a table format. The cells will
	// be rendered with padding.
	GetTableCells func(T) []string

	// function to be called after re-rendering list. Can be nil
	OnRerender func() error

	// set this to true if you don't want to allow manual filtering via '/'
	DisableFilter bool

	// This can be nil if you want to always show the panel
	Hide func() bool

	// Multi-select state: maps item ID to selection status
	selectedItems map[string]bool
	// Function to get ID from item (required for multi-select)
	GetItemID func(T) string
}

var _ ISideListPanel = &SideListPanel[int]{}

type IGui interface {
	HandleClick(v *gocui.View, itemCount int, selectedLine *int, handleSelect func() error) error
	NewSimpleRenderStringTask(getContent func() string) tasks.TaskFunc
	FocusY(selectedLine int, itemCount int, view *gocui.View)
	ShouldRefresh(contextKey string) bool
	GetMainView() *gocui.View
	IsCurrentView(*gocui.View) bool
	FilterString(view *gocui.View) string
	IgnoreStrings() []string
	Update(func() error)

	QueueTask(f func(ctx context.Context)) error
}

func (self *SideListPanel[T]) HandleClick() error {
	itemCount := self.List.Len()
	handleSelect := self.HandleSelect
	selectedLine := &self.SelectedIdx

	if err := self.Gui.HandleClick(self.View, itemCount, selectedLine, handleSelect); err != nil {
		return err
	}

	if self.OnClick != nil {
		selectedItem, err := self.GetSelectedItem()
		if err == nil {
			return self.OnClick(selectedItem)
		}
	}

	return nil
}

func (self *SideListPanel[T]) GetView() *gocui.View {
	return self.View
}

func (self *SideListPanel[T]) HandleSelect() error {
	item, err := self.GetSelectedItem()
	if err != nil {
		if err.Error() != self.NoItemsMessage {
			return err
		}

		if self.NoItemsMessage != "" {
			self.Gui.NewSimpleRenderStringTask(func() string { return self.NoItemsMessage })
		}

		return nil
	}

	self.Refocus()

	return self.renderContext(item)
}

func (self *SideListPanel[T]) renderContext(item T) error {
	if self.ContextState == nil {
		return nil
	}

	key := self.ContextState.GetCurrentContextKey(item)
	if !self.Gui.ShouldRefresh(key) {
		return nil
	}

	mainView := self.Gui.GetMainView()
	mainView.Tabs = self.ContextState.GetMainTabTitles()
	mainView.TabIndex = self.ContextState.mainTabIdx

	task := self.ContextState.GetCurrentMainTab().Render(item)

	return self.Gui.QueueTask(task)
}

func (self *SideListPanel[T]) GetSelectedItem() (T, error) {
	var zero T

	item, ok := self.List.TryGet(self.SelectedIdx)
	if !ok {
		// could probably have a better error here
		return zero, errors.New(self.NoItemsMessage)
	}

	return item, nil
}

func (self *SideListPanel[T]) HandleNextLine() error {
	self.SelectNextLine()

	return self.HandleSelect()
}

func (self *SideListPanel[T]) HandlePrevLine() error {
	self.SelectPrevLine()

	return self.HandleSelect()
}

func (self *SideListPanel[T]) HandleNextMainTab() error {
	if self.ContextState == nil {
		return nil
	}

	self.ContextState.HandleNextMainTab()

	return self.HandleSelect()
}

func (self *SideListPanel[T]) HandlePrevMainTab() error {
	if self.ContextState == nil {
		return nil
	}

	self.ContextState.HandlePrevMainTab()

	return self.HandleSelect()
}

func (self *SideListPanel[T]) Refocus() {
	self.Gui.FocusY(self.SelectedIdx, self.List.Len(), self.View)
}

func (self *SideListPanel[T]) SetItems(items []T) {
	self.List.SetItems(items)
	self.FilterAndSort()
}

func (self *SideListPanel[T]) FilterAndSort() {
	filterString := self.Gui.FilterString(self.View)

	self.List.Filter(func(item T, index int) bool {
		if self.Filter != nil && !self.Filter(item) {
			return false
		}

		if lo.SomeBy(self.Gui.IgnoreStrings(), func(ignore string) bool {
			return lo.SomeBy(self.GetTableCells(item), func(searchString string) bool {
				return strings.Contains(searchString, ignore)
			})
		}) {
			return false
		}

		if filterString != "" {
			return lo.SomeBy(self.GetTableCells(item), func(searchString string) bool {
				return strings.Contains(searchString, filterString)
			})
		}

		return true
	})

	self.List.Sort(self.Sort)

	self.clampSelectedLineIdx()
}

func (self *SideListPanel[T]) RerenderList() error {
	self.FilterAndSort()

	self.Gui.Update(func() error {
		self.View.Clear()
		table := lo.Map(self.List.GetItems(), func(item T, index int) []string {
			return self.GetTableCells(item)
		})
		renderedTable, err := utils.RenderTable(table)
		if err != nil {
			return err
		}
		fmt.Fprint(self.View, renderedTable)

		// Always restore cursor position after re-rendering
		// This ensures highlight stays at the correct line even for non-focused panels
		self.Refocus()

		if self.OnRerender != nil {
			if err := self.OnRerender(); err != nil {
				return err
			}
		}

		if self.Gui.IsCurrentView(self.View) {
			return self.HandleSelect()
		}
		return nil
	})

	return nil
}

func (self *SideListPanel[T]) SetMainTabIndex(index int) {
	if self.ContextState == nil {
		return
	}

	self.ContextState.SetMainTabIndex(index)
}

func (self *SideListPanel[T]) IsFilterDisabled() bool {
	return self.DisableFilter
}

func (self *SideListPanel[T]) IsHidden() bool {
	if self.Hide == nil {
		return false
	}

	return self.Hide()
}

// Multi-select methods

// ToggleSelection toggles the selection state of the item at the current index
func (self *SideListPanel[T]) ToggleSelection() error {
	if self.GetItemID == nil {
		return nil // Multi-select not enabled for this panel
	}

	item, err := self.GetSelectedItem()
	if err != nil {
		return nil
	}

	id := self.GetItemID(item)
	if self.selectedItems == nil {
		self.selectedItems = make(map[string]bool)
	}

	self.selectedItems[id] = !self.selectedItems[id]
	return nil
}

// IsItemSelected returns true if the item at the given index is selected
func (self *SideListPanel[T]) IsItemSelected(idx int) bool {
	if self.GetItemID == nil || self.selectedItems == nil {
		return false
	}

	item, ok := self.List.TryGet(idx)
	if !ok {
		return false
	}

	id := self.GetItemID(item)
	return self.selectedItems[id]
}

// GetSelectedCount returns the number of selected items
func (self *SideListPanel[T]) GetSelectedCount() int {
	if self.selectedItems == nil {
		return 0
	}

	count := 0
	for _, selected := range self.selectedItems {
		if selected {
			count++
		}
	}
	return count
}

// GetSelectedItems returns all selected items
func (self *SideListPanel[T]) GetSelectedItems() []T {
	if self.selectedItems == nil || self.GetItemID == nil {
		return nil
	}

	var selected []T
	for i := 0; i < self.List.Len(); i++ {
		if self.IsItemSelected(i) {
			if item, ok := self.List.TryGet(i); ok {
				selected = append(selected, item)
			}
		}
	}
	return selected
}

// ClearSelection removes all selections
func (self *SideListPanel[T]) ClearSelection() {
	self.selectedItems = make(map[string]bool)
}
