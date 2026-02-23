package panels

import (
	"testing"
)

// makeSidePanel creates a SideListPanel with multi-select enabled.
// Gui and View are intentionally nil — the tested methods don't touch them.
func makeSidePanel(items []string) *SideListPanel[string] {
	fl := NewFilteredList[string]()
	fl.SetItems(items)
	return &SideListPanel[string]{
		ListPanel: ListPanel[string]{
			List:        fl,
			SelectedIdx: 0,
		},
		GetItemID: func(s string) string { return s },
	}
}

func TestSideListPanel_ToggleSelection(t *testing.T) {
	p := makeSidePanel([]string{"alpha", "beta", "gamma"})
	p.SelectedIdx = 0

	// Initially not selected
	if p.IsItemSelected(0) {
		t.Error("item 0 should not be selected initially")
	}

	// Toggle on
	if err := p.ToggleSelection(); err != nil {
		t.Fatalf("ToggleSelection error: %v", err)
	}
	if !p.IsItemSelected(0) {
		t.Error("item 0 should be selected after toggle")
	}

	// Toggle off
	if err := p.ToggleSelection(); err != nil {
		t.Fatalf("ToggleSelection error: %v", err)
	}
	if p.IsItemSelected(0) {
		t.Error("item 0 should be deselected after second toggle")
	}
}

func TestSideListPanel_GetSelectedCount(t *testing.T) {
	p := makeSidePanel([]string{"a", "b", "c"})

	if count := p.GetSelectedCount(); count != 0 {
		t.Errorf("initial GetSelectedCount() = %d, want 0", count)
	}

	p.SelectedIdx = 0
	p.ToggleSelection()
	if count := p.GetSelectedCount(); count != 1 {
		t.Errorf("after 1 selection GetSelectedCount() = %d, want 1", count)
	}

	p.SelectedIdx = 2
	p.ToggleSelection()
	if count := p.GetSelectedCount(); count != 2 {
		t.Errorf("after 2 selections GetSelectedCount() = %d, want 2", count)
	}
}

func TestSideListPanel_GetSelectedItems(t *testing.T) {
	p := makeSidePanel([]string{"alpha", "beta", "gamma"})

	p.SelectedIdx = 0
	p.ToggleSelection()
	p.SelectedIdx = 2
	p.ToggleSelection()

	selected := p.GetSelectedItems()
	if len(selected) != 2 {
		t.Fatalf("GetSelectedItems() len = %d, want 2", len(selected))
	}
	if selected[0] != "alpha" || selected[1] != "gamma" {
		t.Errorf("GetSelectedItems() = %v, want [alpha gamma]", selected)
	}
}

func TestSideListPanel_ClearSelection(t *testing.T) {
	p := makeSidePanel([]string{"a", "b", "c"})

	p.SelectedIdx = 0
	p.ToggleSelection()
	p.SelectedIdx = 1
	p.ToggleSelection()

	if p.GetSelectedCount() != 2 {
		t.Fatalf("expected 2 selected before clear")
	}

	p.ClearSelection()

	if count := p.GetSelectedCount(); count != 0 {
		t.Errorf("after ClearSelection count = %d, want 0", count)
	}
	if items := p.GetSelectedItems(); len(items) != 0 {
		t.Errorf("after ClearSelection GetSelectedItems() = %v, want empty", items)
	}
}

func TestSideListPanel_NoGetItemID_NoOp(t *testing.T) {
	// Without GetItemID, multi-select is disabled — ToggleSelection is a no-op
	fl := NewFilteredList[string]()
	fl.SetItems([]string{"a", "b"})
	p := &SideListPanel[string]{
		ListPanel: ListPanel[string]{
			List:        fl,
			SelectedIdx: 0,
		},
		// GetItemID intentionally nil
	}

	if err := p.ToggleSelection(); err != nil {
		t.Errorf("ToggleSelection with nil GetItemID should return nil error, got %v", err)
	}
	if p.GetSelectedCount() != 0 {
		t.Error("GetSelectedCount should be 0 when GetItemID is nil")
	}
}
