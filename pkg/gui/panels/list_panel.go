package panels

import (
	"github.com/jesseduffield/gocui"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

type ListPanel[T comparable] struct {
	SelectedIdx int
	List        *FilteredList[T]
	View        *gocui.View
}

func (self *ListPanel[T]) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if self.List.Len() > 0 {
		clampedValue = utils.Clamp(value, 0, self.List.Len()-1)
	}

	self.SelectedIdx = clampedValue
}

func (self *ListPanel[T]) clampSelectedLineIdx() {
	clamped := utils.Clamp(self.SelectedIdx, 0, self.List.Len()-1)

	if clamped != self.SelectedIdx {
		self.SelectedIdx = clamped
	}
}

// moves the cursor up or down by the given amount (up for negative values)
func (self *ListPanel[T]) moveSelectedLine(delta int) {
	self.SetSelectedLineIdx(self.SelectedIdx + delta)
}

func (self *ListPanel[T]) SelectNextLine() {
	self.moveSelectedLine(1)
}

func (self *ListPanel[T]) SelectPrevLine() {
	self.moveSelectedLine(-1)
}
