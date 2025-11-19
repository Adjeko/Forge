package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap defines the keybindings for the application
type KeyMap struct {
	Quit      key.Binding
	GitStatus key.Binding
}

// Keys holds the actual key bindings
var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
	GitStatus: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Git Status"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GitStatus, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.GitStatus, k.Quit}}
}

var (
	footerStyle = lipgloss.NewStyle().
		Padding(0, 1)
)

// RenderFooter renders the footer with the given width and help model
func RenderFooter(width int, helpModel help.Model) string {
	if width == 0 {
		return ""
	}

	return footerStyle.Width(width).Render(helpModel.View(Keys))
}
