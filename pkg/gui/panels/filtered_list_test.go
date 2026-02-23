package panels

import (
	"testing"
)

func TestFilteredList_SetAndGet(t *testing.T) {
	fl := NewFilteredList[string]()
	fl.SetItems([]string{"a", "b", "c"})

	if fl.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", fl.Len())
	}

	items := fl.GetItems()
	if len(items) != 3 || items[0] != "a" || items[1] != "b" || items[2] != "c" {
		t.Errorf("GetItems() = %v, want [a b c]", items)
	}
}

func TestFilteredList_Filter(t *testing.T) {
	fl := NewFilteredList[string]()
	fl.SetItems([]string{"apple", "banana", "cherry"})

	// Filter to items containing "a"
	fl.Filter(func(s string, _ int) bool {
		return s == "apple" || s == "banana"
	})

	if fl.Len() != 2 {
		t.Fatalf("after filter Len() = %d, want 2", fl.Len())
	}
	items := fl.GetItems()
	if items[0] != "apple" || items[1] != "banana" {
		t.Errorf("filtered items = %v, want [apple banana]", items)
	}
}

func TestFilteredList_GetAllItems_IgnoresFilter(t *testing.T) {
	fl := NewFilteredList[string]()
	fl.SetItems([]string{"a", "b", "c"})
	fl.Filter(func(s string, _ int) bool { return s == "a" })

	if fl.Len() != 1 {
		t.Fatalf("filtered Len() = %d, want 1", fl.Len())
	}
	all := fl.GetAllItems()
	if len(all) != 3 {
		t.Errorf("GetAllItems() len = %d, want 3 (ignores filter)", len(all))
	}
}

func TestFilteredList_Sort(t *testing.T) {
	fl := NewFilteredList[string]()
	fl.SetItems([]string{"cherry", "apple", "banana"})
	fl.Sort(func(a, b string) bool { return a < b })

	items := fl.GetItems()
	if items[0] != "apple" || items[1] != "banana" || items[2] != "cherry" {
		t.Errorf("sorted items = %v, want [apple banana cherry]", items)
	}
}

func TestFilteredList_TryGet(t *testing.T) {
	fl := NewFilteredList[string]()
	fl.SetItems([]string{"x", "y", "z"})

	item, ok := fl.TryGet(1)
	if !ok || item != "y" {
		t.Errorf("TryGet(1) = (%q, %v), want ('y', true)", item, ok)
	}

	_, ok = fl.TryGet(-1)
	if ok {
		t.Error("TryGet(-1) should return false")
	}

	_, ok = fl.TryGet(99)
	if ok {
		t.Error("TryGet(99) should return false")
	}
}

func TestFilteredList_GetIndex(t *testing.T) {
	fl := NewFilteredList[string]()
	fl.SetItems([]string{"a", "b", "c"})

	if idx := fl.GetIndex("b"); idx != 1 {
		t.Errorf("GetIndex('b') = %d, want 1", idx)
	}
	if idx := fl.GetIndex("z"); idx != -1 {
		t.Errorf("GetIndex('z') = %d, want -1 (not found)", idx)
	}
}
