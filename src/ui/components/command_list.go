package components

import (
	"forge/src/exec"
	"forge/src/output"
	"forge/src/ui/zones"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CommandList component minimal
type CommandList struct {
	Commands []exec.PrimitiveCommand
	Selected int
	Buffer   *output.OutputBuffer
}

func NewCommandList(cmds []exec.PrimitiveCommand, buf *output.OutputBuffer) *CommandList {
	return &CommandList{Commands: cmds, Buffer: buf}
}

func (c *CommandList) View() string {
	var b strings.Builder
	for i, cmd := range c.Commands {
		style := lipgloss.NewStyle()
		if i == c.Selected {
			style = style.Bold(true)
		}
		zones.RegisterZone("command:" + cmd.Cmd)
		b.WriteString(style.Render(cmd.Label) + "\n")
	}
	return b.String()
}
