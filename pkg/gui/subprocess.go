package gui

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/fatih/color"
	"github.com/warl0ck/lazycontainer/pkg/utils"
)

func (gui *Gui) runSubprocess(cmd *exec.Cmd) error {
	return gui.runSubprocessWithMessage(cmd, "")
}

func (gui *Gui) runSubprocessWithMessage(cmd *exec.Cmd, msg string) error {
	// Suspend the GUI
	if err := gui.g.Suspend(); err != nil {
		return fmt.Errorf("failed to suspend GUI: %w", err)
	}

	// Pause background refresh tasks
	gui.PauseBackgroundTasks = true

	// Run the command in foreground
	gui.runCommand(cmd, msg)

	// Resume the GUI
	if err := gui.g.Resume(); err != nil {
		return fmt.Errorf("failed to resume GUI: %w", err)
	}

	// Resume background refresh tasks
	gui.PauseBackgroundTasks = false

	return nil
}

func (gui *Gui) runCommand(cmd *exec.Cmd, msg string) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	stop := make(chan os.Signal, 1)
	defer signal.Stop(stop)

	go func() {
		signal.Notify(stop, os.Interrupt)
		<-stop

		if err := cmd.Process.Kill(); err != nil {
			gui.Log.Error(err)
		}
	}()

	fmt.Fprintf(os.Stdout, "\n%s\n\n", utils.ColoredString("+ "+strings.Join(cmd.Args, " "), color.FgBlue))
	if msg != "" {
		fmt.Fprintf(os.Stdout, "\n%s\n\n", utils.ColoredString(msg, color.FgGreen))
	}
	if err := cmd.Run(); err != nil {
		gui.Log.Error(err)
	}

	cmd.Stdin = nil
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	gui.promptToReturn()
}

func (gui *Gui) promptToReturn() {
	fmt.Fprintf(os.Stdout, "\n%s", utils.ColoredString("Press enter to return...", color.FgGreen))
	_, _ = fmt.Scanln()
}
