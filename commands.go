package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// CommandDefinition defines the structure for a command
type CommandDefinition struct {
	Name      string
	Exe       string
	Args      []string
	Directory string
}

// CommandResult represents the output of a command execution
type CommandResult struct {
	Output string
	Err    error
}

// Execute runs the command and returns the result
func (c CommandDefinition) Execute() (string, error) {
	cmd := exec.Command(c.Exe, c.Args...)
	if c.Directory != "" {
		cmd.Dir = c.Directory
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// RunCommandMsg is a tea.Msg that carries the result of a command
type RunCommandMsg struct {
	Output string
	Err    error
}

// CreateGitStatusCommand creates the specific git status command
func CreateGitStatusCommand() CommandDefinition {
	return CommandDefinition{
		Name:      "Git Status",
		Exe:       "git",
		Args:      []string{"status"},
		Directory: "C:/ADO/CS2",
	}
}

// RunCommand is a tea.Cmd that executes a CommandDefinition
func RunCommand(cmdDef CommandDefinition) tea.Cmd {
	return func() tea.Msg {
		output, err := cmdDef.Execute()
		return RunCommandMsg{Output: output, Err: err}
	}
}
