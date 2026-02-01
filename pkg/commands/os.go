package commands

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type OSCommand struct {
	Log *logrus.Entry
}

func NewOSCommand(log *logrus.Entry) *OSCommand {
	return &OSCommand{
		Log: log,
	}
}

// RunCommand runs a command and returns its output
func (c *OSCommand) RunCommand(cmdStr string) (string, error) {
	c.Log.WithField("command", cmdStr).Debug("Running command")

	args := strings.Fields(cmdStr)
	if len(args) == 0 {
		return "", nil
	}

	cmd := exec.Command(args[0], args[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"command": cmdStr,
			"stderr":  stderr.String(),
			"error":   err,
		}).Debug("Command failed")
		return stderr.String(), err
	}

	return stdout.String(), nil
}

// RunCommandWithArgs runs a command with explicit arguments
func (c *OSCommand) RunCommandWithArgs(name string, args ...string) (string, error) {
	c.Log.WithFields(logrus.Fields{
		"command": name,
		"args":    args,
	}).Debug("Running command with args")

	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"command": name,
			"args":    args,
			"stderr":  stderr.String(),
			"error":   err,
		}).Debug("Command failed")
		return stderr.String(), err
	}

	return stdout.String(), nil
}

// RunCommandWithInput runs a command with stdin input
func (c *OSCommand) RunCommandWithInput(cmdStr string, input string) (string, error) {
	c.Log.WithField("command", cmdStr).Debug("Running command with input")

	args := strings.Fields(cmdStr)
	if len(args) == 0 {
		return "", nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}

	return stdout.String(), nil
}
