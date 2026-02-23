package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/warl0ck/lazycontainer/pkg/config"
)

// ContainerCommand wraps the Apple container CLI
type ContainerCommand struct {
	Log       *logrus.Entry
	OSCommand *OSCommand
	Config    *config.AppConfig
}

// NewContainerCommand creates a new ContainerCommand
func NewContainerCommand(log *logrus.Entry, osCommand *OSCommand, config *config.AppConfig) *ContainerCommand {
	return &ContainerCommand{
		Log:       log,
		OSCommand: osCommand,
		Config:    config,
	}
}

// ListContainers returns all containers
func (c *ContainerCommand) ListContainers(all bool) ([]Container, error) {
	args := []string{"list", "--format", "json"}
	if all {
		args = append(args, "--all")
	}

	output, err := c.OSCommand.RunCommandWithArgs("container", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	if strings.TrimSpace(output) == "" || strings.TrimSpace(output) == "[]" {
		return []Container{}, nil
	}

	var containers []Container
	if err := json.Unmarshal([]byte(output), &containers); err != nil {
		return nil, fmt.Errorf("failed to parse container list: %w", err)
	}

	return containers, nil
}

// InspectContainer returns detailed information about a container
func (c *ContainerCommand) InspectContainer(id string) (*Container, error) {
	output, err := c.OSCommand.RunCommandWithArgs("container", "inspect", "--format", "json", id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container %s: %w", id, err)
	}

	var containers []Container
	if err := json.Unmarshal([]byte(output), &containers); err != nil {
		return nil, fmt.Errorf("failed to parse container inspect: %w", err)
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("container %s not found", id)
	}

	return &containers[0], nil
}

// StartContainer starts a container
func (c *ContainerCommand) StartContainer(id string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "start", id)
	if err != nil {
		return fmt.Errorf("failed to start container %s: %w", id, err)
	}
	return nil
}

// StopContainer stops a container
func (c *ContainerCommand) StopContainer(id string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "stop", id)
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %w", id, err)
	}
	return nil
}

// DeleteContainer removes a container
func (c *ContainerCommand) DeleteContainer(id string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "delete", id)
	if err != nil {
		return fmt.Errorf("failed to delete container %s: %w", id, err)
	}
	return nil
}

// GetContainerLogs fetches container logs
func (c *ContainerCommand) GetContainerLogs(id string, tail int) (string, error) {
	args := []string{"logs", id}
	if tail > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", tail))
	}

	output, err := c.OSCommand.RunCommandWithArgs("container", args...)
	if err != nil {
		return "", fmt.Errorf("failed to get logs for container %s: %w", id, err)
	}

	return output, nil
}

// GetContainerLogsStream returns a command to stream container logs
func (c *ContainerCommand) GetContainerLogsStream(id string, tail int) *exec.Cmd {
	args := []string{"logs", id, "-f"}
	if tail > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", tail))
	}

	return c.OSCommand.BuildCommand("container", args...)
}

// GetContainerExec returns a command to execute an interactive shell in a container
func (c *ContainerCommand) GetContainerExec(id string, cmd string) *exec.Cmd {
	if cmd == "" {
		cmd = "/bin/sh"
	}

	args := []string{"exec", id, "--", cmd}
	return c.OSCommand.BuildCommand("container", args...)
}

// KillContainer sends a signal to a container
func (c *ContainerCommand) KillContainer(id string, signal string) error {
	args := []string{"kill", id}
	if signal != "" {
		args = append(args, "--signal", signal)
	}

	_, err := c.OSCommand.RunCommandWithArgs("container", args...)
	if err != nil {
		return fmt.Errorf("failed to kill container %s: %w", id, err)
	}
	return nil
}

// ListImages returns all images
func (c *ContainerCommand) ListImages() ([]Image, error) {
	output, err := c.OSCommand.RunCommandWithArgs("container", "image", "list", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	if strings.TrimSpace(output) == "" || strings.TrimSpace(output) == "[]" {
		return []Image{}, nil
	}

	var images []Image
	if err := json.Unmarshal([]byte(output), &images); err != nil {
		return nil, fmt.Errorf("failed to parse image list: %w", err)
	}

	return images, nil
}

// PullImage pulls an image from a registry
func (c *ContainerCommand) PullImage(ref string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "image", "pull", ref)
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", ref, err)
	}
	return nil
}

// DeleteImage removes an image
func (c *ContainerCommand) DeleteImage(ref string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "image", "delete", ref)
	if err != nil {
		return fmt.Errorf("failed to delete image %s: %w", ref, err)
	}
	return nil
}

// ListVolumes returns all volumes
func (c *ContainerCommand) ListVolumes() ([]Volume, error) {
	output, err := c.OSCommand.RunCommandWithArgs("container", "volume", "list", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	if strings.TrimSpace(output) == "" || strings.TrimSpace(output) == "[]" {
		return []Volume{}, nil
	}

	var volumes []Volume
	if err := json.Unmarshal([]byte(output), &volumes); err != nil {
		return nil, fmt.Errorf("failed to parse volume list: %w", err)
	}

	return volumes, nil
}

// CreateVolume creates a new volume
func (c *ContainerCommand) CreateVolume(name string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "volume", "create", name)
	if err != nil {
		return fmt.Errorf("failed to create volume %s: %w", name, err)
	}
	return nil
}

// DeleteVolume removes a volume
func (c *ContainerCommand) DeleteVolume(name string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "volume", "delete", name)
	if err != nil {
		return fmt.Errorf("failed to delete volume %s: %w", name, err)
	}
	return nil
}

// ListNetworks returns all networks
func (c *ContainerCommand) ListNetworks() ([]Network, error) {
	output, err := c.OSCommand.RunCommandWithArgs("container", "network", "list", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	if strings.TrimSpace(output) == "" || strings.TrimSpace(output) == "[]" {
		return []Network{}, nil
	}

	var networks []Network
	if err := json.Unmarshal([]byte(output), &networks); err != nil {
		return nil, fmt.Errorf("failed to parse network list: %w", err)
	}

	return networks, nil
}

// CreateNetwork creates a new network
func (c *ContainerCommand) CreateNetwork(name string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "network", "create", name)
	if err != nil {
		return fmt.Errorf("failed to create network %s: %w", name, err)
	}
	return nil
}

// DeleteNetwork removes a network
func (c *ContainerCommand) DeleteNetwork(name string) error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "network", "delete", name)
	if err != nil {
		return fmt.Errorf("failed to delete network %s: %w", name, err)
	}
	return nil
}

// SystemStatus returns the status of container services
func (c *ContainerCommand) SystemStatus() (string, error) {
	output, err := c.OSCommand.RunCommandWithArgs("container", "system", "status")
	if err != nil {
		return "", fmt.Errorf("failed to get system status: %w", err)
	}
	return output, nil
}

// SystemStart starts container services
func (c *ContainerCommand) SystemStart() error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "system", "start")
	if err != nil {
		return fmt.Errorf("failed to start container services: %w", err)
	}
	return nil
}

// SystemStop stops all container services
func (c *ContainerCommand) SystemStop() error {
	_, err := c.OSCommand.RunCommandWithArgs("container", "system", "stop")
	if err != nil {
		return fmt.Errorf("failed to stop container services: %w", err)
	}
	return nil
}

// GetContainerStats returns resource usage stats
func (c *ContainerCommand) GetContainerStats(ids ...string) (map[string]ContainerStats, error) {
	args := []string{"stats", "--no-stream", "--format", "json"}
	args = append(args, ids...)

	_, err := c.OSCommand.RunCommandWithArgs("container", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}

	// Parse stats output - format may vary
	stats := make(map[string]ContainerStats)

	// TODO: Parse the JSON output once we see the actual format
	// from `container stats --format json`

	return stats, nil
}
