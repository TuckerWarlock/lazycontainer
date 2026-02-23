package panels

import (
	"testing"
)

// makeListPanel creates a ListPanel with the given items and initial SelectedIdx=0.
// Note: View is nil intentionally — navigation methods don't touch it.
func makeListPanel(items []string) *ListPanel[string] {
	fl := NewFilteredList[string]()
	fl.SetItems(items)
	return &ListPanel[string]{
		List:        fl,
		SelectedIdx: 0,
	}
}

func TestListPanel_SelectNextLine(t *testing.T) {
	lp := makeListPanel([]string{"a", "b", "c"})

	lp.SelectNextLine()
	if lp.SelectedIdx != 1 {
		t.Errorf("after SelectNextLine SelectedIdx = %d, want 1", lp.SelectedIdx)
	}

	lp.SelectNextLine()
	if lp.SelectedIdx != 2 {
		t.Errorf("after 2nd SelectNextLine SelectedIdx = %d, want 2", lp.SelectedIdx)
	}

	// At end: should clamp to last index
	lp.SelectNextLine()
	if lp.SelectedIdx != 2 {
		t.Errorf("SelectNextLine at end SelectedIdx = %d, want 2 (clamped)", lp.SelectedIdx)
	}
}

func TestListPanel_SelectPrevLine(t *testing.T) {
	lp := makeListPanel([]string{"a", "b", "c"})
	lp.SelectedIdx = 2

	lp.SelectPrevLine()
	if lp.SelectedIdx != 1 {
		t.Errorf("after SelectPrevLine SelectedIdx = %d, want 1", lp.SelectedIdx)
	}

	lp.SelectPrevLine()
	if lp.SelectedIdx != 0 {
		t.Errorf("after 2nd SelectPrevLine SelectedIdx = %d, want 0", lp.SelectedIdx)
	}

	// At start: should clamp to 0
	lp.SelectPrevLine()
	if lp.SelectedIdx != 0 {
		t.Errorf("SelectPrevLine at start SelectedIdx = %d, want 0 (clamped)", lp.SelectedIdx)
	}
}

func TestListPanel_SetSelectedLineIdx(t *testing.T) {
	lp := makeListPanel([]string{"a", "b", "c"}) // indices 0-2

	lp.SetSelectedLineIdx(1)
	if lp.SelectedIdx != 1 {
		t.Errorf("SetSelectedLineIdx(1) = %d, want 1", lp.SelectedIdx)
	}

	// Clamp below 0
	lp.SetSelectedLineIdx(-5)
	if lp.SelectedIdx != 0 {
		t.Errorf("SetSelectedLineIdx(-5) = %d, want 0 (clamped)", lp.SelectedIdx)
	}

	// Clamp above max
	lp.SetSelectedLineIdx(100)
	if lp.SelectedIdx != 2 {
		t.Errorf("SetSelectedLineIdx(100) = %d, want 2 (clamped to len-1)", lp.SelectedIdx)
	}
}
