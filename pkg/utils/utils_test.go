package utils

import (
	"runtime"
	"strings"
	"testing"
)

func TestClamp(t *testing.T) {
	if got := Clamp(5, 1, 10); got != 5 {
		t.Errorf("Clamp(5,1,10) = %d, want 5", got)
	}
	if got := Clamp(0, 1, 10); got != 1 {
		t.Errorf("Clamp(0,1,10) = %d, want 1 (below min)", got)
	}
	if got := Clamp(15, 1, 10); got != 10 {
		t.Errorf("Clamp(15,1,10) = %d, want 10 (above max)", got)
	}
	if got := Clamp(1, 1, 10); got != 1 {
		t.Errorf("Clamp(1,1,10) = %d, want 1 (at min boundary)", got)
	}
	if got := Clamp(10, 1, 10); got != 10 {
		t.Errorf("Clamp(10,1,10) = %d, want 10 (at max boundary)", got)
	}
}

func TestMax(t *testing.T) {
	if got := Max(3, 5); got != 5 {
		t.Errorf("Max(3,5) = %d, want 5", got)
	}
	if got := Max(5, 3); got != 5 {
		t.Errorf("Max(5,3) = %d, want 5", got)
	}
	if got := Max(4, 4); got != 4 {
		t.Errorf("Max(4,4) = %d, want 4", got)
	}
}

func TestSplitLines(t *testing.T) {
	if got := SplitLines(""); len(got) != 0 {
		t.Errorf("SplitLines('') = %v, want empty slice", got)
	}
	if got := SplitLines("hello"); len(got) != 1 || got[0] != "hello" {
		t.Errorf("SplitLines single line = %v", got)
	}
	got := SplitLines("a\nb\nc")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("SplitLines multiline = %v", got)
	}
	// Windows CRLF
	got = SplitLines("a\r\nb")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("SplitLines CRLF = %v", got)
	}
}

func TestFormatBinaryBytes(t *testing.T) {
	if got := FormatBinaryBytes(0); got != "0B" {
		t.Errorf("FormatBinaryBytes(0) = %q, want '0B'", got)
	}
	if got := FormatBinaryBytes(512); !strings.HasSuffix(got, "B") {
		t.Errorf("FormatBinaryBytes(512) = %q, want bytes", got)
	}
	// The function uses strict > comparison so exactly 1024 stays in bytes unit
	if got := FormatBinaryBytes(1024); got != "1024.00B" {
		t.Errorf("FormatBinaryBytes(1024) = %q, want '1024.00B'", got)
	}
	// 1025 is the first value that crosses the 1024 threshold and converts to kiB
	if got := FormatBinaryBytes(1025); got != "1.00kiB" {
		t.Errorf("FormatBinaryBytes(1025) = %q, want '1.00kiB'", got)
	}
	// 1024*1024 stays in kiB since it is not > 1024*1024
	if got := FormatBinaryBytes(1024 * 1024); got != "1024.00kiB" {
		t.Errorf("FormatBinaryBytes(1MiB) = %q, want '1024.00kiB'", got)
	}
	// 1024*1024+1 crosses into MiB
	if got := FormatBinaryBytes(1024*1024 + 1); got != "1.00MiB" {
		t.Errorf("FormatBinaryBytes(1MiB+1) = %q, want '1.00MiB'", got)
	}
}

func TestFormatDecimalBytes(t *testing.T) {
	if got := FormatDecimalBytes(0); got != "0B" {
		t.Errorf("FormatDecimalBytes(0) = %q, want '0B'", got)
	}
	if got := FormatDecimalBytes(500); !strings.HasSuffix(got, "B") {
		t.Errorf("FormatDecimalBytes(500) = %q, want bytes", got)
	}
	// The function uses strict > comparison so exactly 1000 stays in bytes unit
	if got := FormatDecimalBytes(1000); got != "1000.00B" {
		t.Errorf("FormatDecimalBytes(1000) = %q, want '1000.00B'", got)
	}
	// 1001 is the first value that crosses the 1000 threshold and converts to kB
	if got := FormatDecimalBytes(1001); got != "1.00kB" {
		t.Errorf("FormatDecimalBytes(1001) = %q, want '1.00kB'", got)
	}
	// 1000*1000 stays in kB since it is not > 1000*1000
	if got := FormatDecimalBytes(1000 * 1000); got != "1000.00kB" {
		t.Errorf("FormatDecimalBytes(1MB) = %q, want '1000.00kB'", got)
	}
	// 1000*1000+1 crosses into MB
	if got := FormatDecimalBytes(1000*1000 + 1); got != "1.00MB" {
		t.Errorf("FormatDecimalBytes(1MB+1) = %q, want '1.00MB'", got)
	}
}

func TestDecolorise(t *testing.T) {
	// Plain string passes through unchanged
	if got := Decolorise("hello"); got != "hello" {
		t.Errorf("Decolorise plain = %q, want 'hello'", got)
	}
	// Strips ANSI color codes
	colored := "\x1B[32mgreen\x1B[0m"
	if got := Decolorise(colored); got != "green" {
		t.Errorf("Decolorise colored = %q, want 'green'", got)
	}
}

func TestWithPadding(t *testing.T) {
	// Pads short string
	got := WithPadding("hi", 5)
	if len(got) != 5 {
		t.Errorf("WithPadding('hi', 5) len = %d, want 5", len(got))
	}
	if got[:2] != "hi" {
		t.Errorf("WithPadding('hi', 5) = %q, should start with 'hi'", got)
	}
	// Does not truncate longer string
	got = WithPadding("hello world", 5)
	if got != "hello world" {
		t.Errorf("WithPadding('hello world', 5) = %q, want unchanged", got)
	}
}

func TestRenderTable(t *testing.T) {
	// Empty table
	result, err := RenderTable([][]string{})
	if err != nil || result != "" {
		t.Errorf("RenderTable empty: got %q, %v", result, err)
	}
	// Single row
	result, err = RenderTable([][]string{{"a", "b", "c"}})
	if err != nil {
		t.Fatalf("RenderTable single row error: %v", err)
	}
	if !strings.Contains(result, "a") || !strings.Contains(result, "b") || !strings.Contains(result, "c") {
		t.Errorf("RenderTable single row = %q, missing expected values", result)
	}
	// Multi-row alignment: all rows must have same column count
	result, err = RenderTable([][]string{
		{"name", "status"},
		{"alice", "running"},
	})
	if err != nil {
		t.Fatalf("RenderTable multi-row error: %v", err)
	}
	if !strings.Contains(result, "name") || !strings.Contains(result, "alice") {
		t.Errorf("RenderTable multi-row = %q, missing expected values", result)
	}
	// Mismatched row lengths returns error
	_, err = RenderTable([][]string{
		{"a", "b"},
		{"x"},
	})
	if err == nil {
		t.Error("RenderTable mismatched rows should return error")
	}
}

func TestCopyToClipboard(t *testing.T) {
	// Skip on non-Darwin systems since pbcopy is macOS-only
	if runtime.GOOS != "darwin" {
		t.Skip("pbcopy not available on non-macOS systems")
	}

	testString := "test-container-id"
	err := CopyToClipboard(testString)
	if err != nil {
		t.Errorf("CopyToClipboard failed: %v", err)
	}
}
