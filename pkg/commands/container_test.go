package commands

import (
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestContainerJSONParsing(t *testing.T) {
	jsonData := `[{"status":"running","configuration":{"dns":{"searchDomains":[],"options":[],"nameservers":[]},"runtimeHandler":"container-runtime-linux","labels":{},"publishedPorts":[{"count":1,"hostAddress":"0.0.0.0","proto":"tcp","containerPort":80,"hostPort":8080}],"image":{"reference":"docker.io/library/alpine:latest","descriptor":{"size":9218,"mediaType":"application/vnd.oci.image.index.v1+json","digest":"sha256:abc123"}},"id":"test-container","networks":[{"options":{"hostname":"test"},"network":"default"}],"readOnly":false,"virtualization":false,"ssh":false,"sysctls":{},"initProcess":{"supplementalGroups":[],"user":{"id":{"gid":0,"uid":0}},"executable":"sleep","workingDirectory":"/","terminal":false,"arguments":["3600"],"rlimits":[],"environment":["PATH=/usr/local/bin"]},"publishedSockets":[],"rosetta":false,"resources":{"memoryInBytes":1073741824,"cpus":4},"platform":{"os":"linux","architecture":"arm64"},"mounts":[]},"startedDate":791617985.958395,"networks":[{"hostname":"test","network":"default","ipv4Address":"192.168.64.8/24","ipv6Address":"fd::1/64","macAddress":"aa:bb:cc:dd:ee:ff","ipv4Gateway":"192.168.64.1"}]}]`

	var containers []Container
	if err := json.Unmarshal([]byte(jsonData), &containers); err != nil {
		t.Fatalf("Failed to parse container JSON: %v", err)
	}

	if len(containers) != 1 {
		t.Fatalf("Expected 1 container, got %d", len(containers))
	}

	c := containers[0]

	if c.GetID() != "test-container" {
		t.Errorf("Expected ID 'test-container', got '%s'", c.GetID())
	}

	if c.GetStatus() != "running" {
		t.Errorf("Expected status 'running', got '%s'", c.GetStatus())
	}

	if !c.IsRunning() {
		t.Errorf("Expected IsRunning() to be true")
	}

	if c.GetImage() != "docker.io/library/alpine:latest" {
		t.Errorf("Expected image 'docker.io/library/alpine:latest', got '%s'", c.GetImage())
	}

	if c.GetCPUs() != 4 {
		t.Errorf("Expected 4 CPUs, got %d", c.GetCPUs())
	}

	if c.GetPorts() != "8080:80" {
		t.Errorf("Expected ports '8080:80', got '%s'", c.GetPorts())
	}

	if c.GetIPAddress() != "192.168.64.8" {
		t.Errorf("Expected IP '192.168.64.8', got '%s'", c.GetIPAddress())
	}

	t.Logf("Container parsed successfully: %+v", c.GetName())
}

func TestImageJSONParsing(t *testing.T) {
	jsonData := `[{"reference":"docker.io/library/alpine:latest","descriptor":{"size":9218,"digest":"sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659","mediaType":"application/vnd.oci.image.index.v1+json"}}]`

	var images []Image
	if err := json.Unmarshal([]byte(jsonData), &images); err != nil {
		t.Fatalf("Failed to parse image JSON: %v", err)
	}

	if len(images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(images))
	}

	i := images[0]

	if i.Reference != "docker.io/library/alpine:latest" {
		t.Errorf("Expected reference 'docker.io/library/alpine:latest', got '%s'", i.Reference)
	}

	if i.GetDigestShort() != "25109184c71b" {
		t.Errorf("Expected digest short '25109184c71b', got '%s'", i.GetDigestShort())
	}

	t.Logf("Image parsed successfully: %s", i.Reference)
}

func TestVolumeJSONParsing(t *testing.T) {
	jsonData := `[{"sizeInBytes":549755813888,"options":{},"source":"/path/to/volume.img","name":"test-volume","format":"ext4","driver":"local","createdAt":791617985.085557,"labels":{}}]`

	var volumes []Volume
	if err := json.Unmarshal([]byte(jsonData), &volumes); err != nil {
		t.Fatalf("Failed to parse volume JSON: %v", err)
	}

	if len(volumes) != 1 {
		t.Fatalf("Expected 1 volume, got %d", len(volumes))
	}

	v := volumes[0]

	if v.Name != "test-volume" {
		t.Errorf("Expected name 'test-volume', got '%s'", v.Name)
	}

	if v.Format != "ext4" {
		t.Errorf("Expected format 'ext4', got '%s'", v.Format)
	}

	t.Logf("Volume parsed successfully: %s", v.Name)
}

func TestNetworkJSONParsing(t *testing.T) {
	jsonData := `[{"state":"running","status":{"ipv4Subnet":"192.168.64.0/24","ipv6Subnet":"fd::0/64","ipv4Gateway":"192.168.64.1"},"id":"default","config":{"mode":"nat","labels":{},"id":"default","creationDate":791602818.713889}}]`

	var networks []Network
	if err := json.Unmarshal([]byte(jsonData), &networks); err != nil {
		t.Fatalf("Failed to parse network JSON: %v", err)
	}

	if len(networks) != 1 {
		t.Fatalf("Expected 1 network, got %d", len(networks))
	}

	n := networks[0]

	if n.ID != "default" {
		t.Errorf("Expected ID 'default', got '%s'", n.ID)
	}

	if n.State != "running" {
		t.Errorf("Expected state 'running', got '%s'", n.State)
	}

	if n.Config.Mode != "nat" {
		t.Errorf("Expected mode 'nat', got '%s'", n.Config.Mode)
	}

	t.Logf("Network parsed successfully: %s", n.ID)
}

func TestGetContainerLogsStream(t *testing.T) {
	log := GetTestLogger()
	osCmd := NewOSCommand(log)
	containerCmd := NewContainerCommand(log, osCmd, nil)

	// Test that the command is properly constructed
	cmd := containerCmd.GetContainerLogsStream("test-container", 100)

	if cmd == nil {
		t.Fatalf("Expected non-nil command")
	}

	// Check that args include logs, container ID, and follow flag
	args := cmd.Args
	if len(args) < 4 {
		t.Fatalf("Expected at least 4 args, got %d: %v", len(args), args)
	}

	// args[0] is the command name, args[1] should be "logs"
	if args[1] != "logs" {
		t.Errorf("Expected args[1] to be 'logs', got '%s'", args[1])
	}

	// args[2] should be the container ID
	if args[2] != "test-container" {
		t.Errorf("Expected args[2] to be 'test-container', got '%s'", args[2])
	}

	// args[3] should be "-f"
	if args[3] != "-f" {
		t.Errorf("Expected args[3] to be '-f', got '%s'", args[3])
	}

	t.Logf("GetContainerLogsStream command constructed correctly: %v", args)
}

func TestGetContainerExec(t *testing.T) {
	log := GetTestLogger()
	osCmd := NewOSCommand(log)
	containerCmd := NewContainerCommand(log, osCmd, nil)

	// Test with default shell
	cmd := containerCmd.GetContainerExec("test-container", "")

	if cmd == nil {
		t.Fatalf("Expected non-nil command")
	}

	args := cmd.Args
	if len(args) < 5 {
		t.Fatalf("Expected at least 5 args, got %d: %v", len(args), args)
	}

	// args[1] should be "exec"
	if args[1] != "exec" {
		t.Errorf("Expected args[1] to be 'exec', got '%s'", args[1])
	}

	// args[2] should be the container ID
	if args[2] != "test-container" {
		t.Errorf("Expected args[2] to be 'test-container', got '%s'", args[2])
	}

	// args[3] should be "--"
	if args[3] != "--" {
		t.Errorf("Expected args[3] to be '--', got '%s'", args[3])
	}

	// args[4] should be the default shell "/bin/sh"
	if args[4] != "/bin/sh" {
		t.Errorf("Expected args[4] to be '/bin/sh', got '%s'", args[4])
	}

	t.Logf("GetContainerExec default shell command correct: %v", args)

	// Test with custom command
	cmdCustom := containerCmd.GetContainerExec("test-container", "ls -la")

	argsCustom := cmdCustom.Args
	if len(argsCustom) < 5 {
		t.Fatalf("Expected at least 5 args for custom command, got %d: %v", len(argsCustom), argsCustom)
	}

	// args[4] should be the custom command
	if argsCustom[4] != "ls -la" {
		t.Errorf("Expected args[4] to be 'ls -la', got '%s'", argsCustom[4])
	}

	t.Logf("GetContainerExec custom command correct: %v", argsCustom)
}

func TestGetContainerStatsStream(t *testing.T) {
	log := GetTestLogger()
	osCmd := NewOSCommand(log)
	containerCmd := NewContainerCommand(log, osCmd, nil)

	cmd := containerCmd.GetContainerStatsStream("test-container")

	if cmd == nil {
		t.Fatalf("Expected non-nil command")
	}

	args := cmd.Args
	if len(args) < 5 {
		t.Fatalf("Expected at least 5 args, got %d: %v", len(args), args)
	}

	// args[1] should be "stats"
	if args[1] != "stats" {
		t.Errorf("Expected args[1] to be 'stats', got '%s'", args[1])
	}

	// args[2] should be the container ID
	if args[2] != "test-container" {
		t.Errorf("Expected args[2] to be 'test-container', got '%s'", args[2])
	}

	// args[3] should be "--format"
	if args[3] != "--format" {
		t.Errorf("Expected args[3] to be '--format', got '%s'", args[3])
	}

	// args[4] should be "json"
	if args[4] != "json" {
		t.Errorf("Expected args[4] to be 'json', got '%s'", args[4])
	}

	t.Logf("GetContainerStatsStream command constructed correctly: %v", args)
}

func GetTestLogger() *logrus.Entry {
	return logrus.WithField("test", true)
}
