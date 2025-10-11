package help

import (
	"forge/src/ui/accessibility"
	"forge/src/ui/zones"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Overlay represents a help overlay listing hotkey actions.
type Overlay struct {
	Visible bool
}

func NewOverlay() *Overlay { return &Overlay{Visible: false} }

func (o *Overlay) Toggle() { o.Visible = !o.Visible }

func (o *Overlay) View() string {
	if o == nil || !o.Visible {
		return ""
	}
	zones.RegisterZone("help:overlay")
	actions := accessibility.List()
	if len(actions) == 0 {
		return lipgloss.NewStyle().Faint(true).Render("No actions registered")
	}
	// Build table lines
	lines := []string{"Hotkeys | Zone | Description"}
	for _, a := range actions {
		keys := strings.Join(a.Keys, "+")
		lines = append(lines, keys+" | "+a.ZoneID+" | "+a.Description)
	}
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Render(strings.Join(lines, "\n"))
	title := lipgloss.NewStyle().Bold(true).Underline(true).Render("Help (? to hide)")
	return lipgloss.JoinVertical(lipgloss.Left, title, box)
}
