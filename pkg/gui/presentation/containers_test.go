package presentation

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/warl0ck/lazycontainer/pkg/commands"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

// parseTestContainer builds a *commands.Container from minimal JSON for testing.
func parseTestContainer(t *testing.T, id, status string) *commands.Container {
	t.Helper()
	jsonStr := `[{"status":"` + status + `","configuration":{"id":"` + id + `","dns":{"searchDomains":[],"options":[],"nameservers":[]},"runtimeHandler":"","labels":{},"publishedPorts":[],"image":{"reference":"docker.io/library/alpine:latest","descriptor":{"size":0,"mediaType":"","digest":""}},"networks":[],"readOnly":false,"virtualization":false,"ssh":false,"sysctls":{},"initProcess":{"supplementalGroups":[],"user":{"id":{"gid":0,"uid":0}},"executable":"sleep","workingDirectory":"/","terminal":false,"arguments":[],"rlimits":[],"environment":[]},"publishedSockets":[],"rosetta":false,"resources":{"memoryInBytes":0,"cpus":0},"platform":{"os":"linux","architecture":"arm64"},"mounts":[]},"startedDate":0,"networks":[]}]`
	var containers []commands.Container
	if err := json.Unmarshal([]byte(jsonStr), &containers); err != nil {
		t.Fatalf("failed to parse container JSON: %v", err)
	}
	return &containers[0]
}

// TestGetContainerDisplayStringsReturns5Columns verifies that the function
// always returns exactly 5 display columns.
func TestGetContainerDisplayStringsReturns5Columns(t *testing.T) {
	c := parseTestContainer(t, "test-container", "running")
	cols := GetContainerDisplayStrings(c)
	if len(cols) != 5 {
		t.Errorf("expected 5 columns, got %d: %v", len(cols), cols)
	}
}

// TestGetContainerDisplayStringsRunningVsStopped verifies that a running
// container and a stopped container produce different values in the status
// column (index 0).
func TestGetContainerDisplayStringsRunningVsStopped(t *testing.T) {
	running := parseTestContainer(t, "running-container", "running")
	stopped := parseTestContainer(t, "stopped-container", "stopped")

	runningCols := GetContainerDisplayStrings(running)
	stoppedCols := GetContainerDisplayStrings(stopped)

	// Decolorise before comparing so ANSI codes do not affect the comparison.
	runningStatus := utils.Decolorise(runningCols[0])
	stoppedStatus := utils.Decolorise(stoppedCols[0])

	if runningStatus == stoppedStatus {
		t.Errorf(
			"expected different status symbols for running vs stopped, but both are %q",
			runningStatus,
		)
	}
}

// TestGetContainerDisplayStringsImagePrefixStripping verifies that the
// "docker.io/library/" prefix is stripped from the image column (index 4).
func TestGetContainerDisplayStringsImagePrefixStripping(t *testing.T) {
	c := parseTestContainer(t, "prefix-test", "running")

	// Confirm the raw image reference still includes the full prefix.
	if !strings.HasPrefix(c.GetImage(), "docker.io/library/") {
		t.Fatalf("test setup error: expected image reference to start with 'docker.io/library/', got %q", c.GetImage())
	}

	cols := GetContainerDisplayStrings(c)
	imageCol := utils.Decolorise(cols[4])

	if strings.Contains(imageCol, "docker.io/library/") {
		t.Errorf("image column should not contain 'docker.io/library/' prefix, got %q", imageCol)
	}

	// The stripped value should contain "alpine:latest".
	if !strings.Contains(imageCol, "alpine:latest") {
		t.Errorf("image column should contain 'alpine:latest', got %q", imageCol)
	}
}
