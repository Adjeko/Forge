package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sewworkspacemanager/internal/domain"
)

// (CommandNode now lives in internal/domain)

// CommandTree ist eine navigierbare Baum-Komponente mit Detailbereich.
type CommandTree struct {
	root    *domain.CommandNode
	flat    []*domain.CommandNode // flattened visible nodes in display order
	cursor  int
	wTree   int
	wDetail int
	help    help.Model
	keys    treeKeyMap
	Title   string
}

// treeKeyMap implementiert help.KeyMap für den Baum-Dialog.
type treeKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Esc    key.Binding
	Close  key.Binding // Strg+A
	Create key.Binding // optional (Strg+N) für Abläufe neu
}

func (k treeKeyMap) ShortHelp() []key.Binding {
	if k.Create.Keys() != nil {
		return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Enter, k.Create, k.Esc}
	}
	return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Enter, k.Esc}
}
func (k treeKeyMap) FullHelp() [][]key.Binding {
	row := []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Enter, k.Esc, k.Close}
	if k.Create.Keys() != nil {
		row = append(row, k.Create)
	}
	return [][]key.Binding{row}
}

func NewCommandTree(root *domain.CommandNode) *CommandTree {
	ct := &CommandTree{root: root, cursor: 0, wTree: 40, wDetail: 20}
	ct.help = help.New()
	// Stil für Hilfe konsistent mit globalen Styles (falls vorhanden)
	ct.help.Styles.ShortKey = hotkeyStyle
	ct.help.Styles.ShortDesc = textStyle
	ct.help.Styles.FullKey = hotkeyStyle
	ct.help.Styles.FullDesc = textStyle
	ct.keys = treeKeyMap{
		Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "Hoch")),
		Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "Runter")),
		Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("←", "Einklappen")),
		Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "Aufklappen")),
		Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "Ausführen")),
		Esc:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("Esc", "Schließen")),
		Close: key.NewBinding(key.WithKeys("ctrl+b"), key.WithHelp("^+B", "Schließen")),
	}
	// Root standardmäßig aufgeklappt
	root.Expanded = true
	ct.Title = "Befehle"
	ct.refresh()
	return ct
}

// NewCommandTreeWithTitle erstellt einen Befehlsbaum mit benutzerdefiniertem Dialogtitel.
func NewCommandTreeWithTitle(root *domain.CommandNode, title string) *CommandTree {
	ct := NewCommandTree(root)
	if title != "" {
		ct.Title = title
	}
	return ct
}

func (c *CommandTree) Width() int { return c.wTree }
func (c *CommandTree) SetWidth(w int) {
	if w > 10 {
		c.wTree = w
	}
}
func (c *CommandTree) SetDetailWidth(w int) {
	if w > 10 {
		c.wDetail = w
	}
}

// SelectedCommand gibt den Knoten zurück, wenn die aktuelle Auswahl ausführbar ist.
func (c *CommandTree) SelectedCommand() *domain.CommandNode {
	if c.cursor >= 0 && c.cursor < len(c.flat) {
		n := c.flat[c.cursor]
		if n.Command != "" {
			return n
		}
	}
	return nil
}

// Update verarbeitet Navigations-Tasten.
func (c *CommandTree) Update(msg tea.Msg) (*CommandTree, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		key := m.String()
		switch key {
		case "up":
			if c.cursor > 0 {
				c.cursor--
			}
		case "down":
			if c.cursor < len(c.flat)-1 {
				c.cursor++
			}
			// einklappen oder zum Elternknoten gehen
			cur := c.current()
			if cur != nil {
				if cur.Expanded && len(cur.Children) > 0 {
					cur.Expanded = false
					c.refresh()
				} else if cur.Parent != nil {
					// zum Elternknoten bewegen
					p := cur.Parent
					for i, n := range c.flat {
						if n == p {
							c.cursor = i
							break
						}
					}
				}
			}
		case "right":
			cur := c.current()
			if cur != nil && len(cur.Children) > 0 {
				if !cur.Expanded {
					cur.Expanded = true
					c.refresh()
				} else if cur.Expanded { // zum ersten Kind springen
					if len(cur.Children) > 0 {
						for i, n := range c.flat {
							if n == cur.Children[0] {
								c.cursor = i
								break
							}
						}
					}
				}
			}
		}
	}
	return c, nil
}

func (c *CommandTree) current() *domain.CommandNode {
	if c.cursor >= 0 && c.cursor < len(c.flat) {
		return c.flat[c.cursor]
	}
	return nil
}

func (c *CommandTree) refresh() {
	// Root-Knoten nicht mehr anzeigen: starte Traversierung direkt bei seinen Kindern.
	c.flat = c.flat[:0]
	var walk func(n *domain.CommandNode, depth int)
	walk = func(n *domain.CommandNode, depth int) {
		c.flat = append(c.flat, n)
		if depth >= 2 { // depth limit beibehalten
			return
		}
		if n.Expanded {
			for _, ch := range n.Children {
				walk(ch, depth+1)
			}
		}
	}
	// Kinder des Root als depth 0 behandeln
	for _, ch := range c.root.Children {
		walk(ch, 0)
	}
	if c.cursor >= len(c.flat) {
		c.cursor = len(c.flat) - 1
	}
	if c.cursor < 0 {
		c.cursor = 0
	}
}

func (c *CommandTree) View() string {
	// Baum-Liste links
	title := c.Title
	if title == "" {
		title = "Befehle"
	}
	header := RenderDialogHeader(title, c.wTree+c.wDetail)
	var b strings.Builder
	for i, n := range c.flat {
		depth := domain.DepthOf(n)
		indent := strings.Repeat("  ", depth)
		marker := " "
		if len(n.Children) > 0 {
			if n.Expanded {
				marker = "▾"
			} else {
				marker = "▸"
			}
		}
		// Icon Logik: Ordner / Befehl / gemischt
		icon := ""
		if len(n.Children) > 0 && n.Command == "" {
			icon = "📁"
		} else if len(n.Children) > 0 && n.Command != "" {
			icon = "🗂"
		} else if n.Command != "" {
			icon = "⚙"
		}
		label := n.Label
		if icon != "" {
			label = icon + "  " + label
		}
		line := indent + marker + " " + label
		// clip
		runes := []rune(line)
		if len(runes) > c.wTree-1 {
			line = string(runes[:c.wTree-1])
		}
		rowStyle := lipgloss.NewStyle()
		if len(n.Children) > 0 && n.Command == "" { // reiner Ordner
			rowStyle = rowStyle.Foreground(lipgloss.Color("#E30018"))
		} else if n.Command != "" {
			rowStyle = rowStyle.Foreground(lipgloss.Color("#CCCCCC"))
		}
		if i == c.cursor {
			rowStyle = rowStyle.Background(lipgloss.Color("#B50013")).Foreground(lipgloss.Color("#FFFFFF"))
		}
		line = rowStyle.Render(domain.PadRight(line, c.wTree-1))
		b.WriteString(line)
		b.WriteString("\n")
	}
	left := b.String()
	// Detailbereich
	detail := c.detailView()
	// Nebeneinander zusammensetzen
	linesL := strings.Split(strings.TrimSuffix(left, "\n"), "\n")
	linesD := strings.Split(detail, "\n")
	max := len(linesL)
	if len(linesD) > max {
		max = len(linesD)
	}
	var out []string
	for i := 0; i < max; i++ {
		var lcol, rcol string
		if i < len(linesL) {
			lcol = domain.PadRight(linesL[i], c.wTree)
		} else {
			lcol = strings.Repeat(" ", c.wTree)
		}
		if i < len(linesD) {
			rcol = domain.PadRight(linesD[i], c.wDetail)
		} else {
			rcol = strings.Repeat(" ", c.wDetail)
		}
		out = append(out, lcol+rcol)
	}
	mainBlock := header + "\n\n" + strings.Join(out, "\n")
	// Hilfe anhängen (wrapt user bereitgestelltes Hilfe-View automatisch)
	helpView := c.help.View(c.keys)
	// Breite sicherstellen; darf gesamte Dialogbreite nutzen
	if helpView != "" {
		mainBlock += "\n" + helpView
	}
	return mainBlock
}

func (c *CommandTree) detailView() string {
	n := c.current()
	if n == nil {
		return ""
	}
	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(n.Label))
	if n.Description != "" {
		lines = append(lines, n.Description)
	}
	if n.Command != "" {
		lines = append(lines, "")
		lines = append(lines, "Befehl: "+n.Command)
		if len(n.Args) > 0 {
			lines = append(lines, "Argumente: "+strings.Join(n.Args, " "))
		}
		if n.WorkDir != "" {
			lines = append(lines, "Verzeichnis: "+n.WorkDir)
		}
	} else if len(n.Children) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Ordner mit "+domain.Itoa(len(n.Children))+" Einträgen")
	}
	return strings.Join(lines, "\n")
}

// (Hilfsfunktionen nach internal/domain verschoben)

// legendLines entfernt: help.Model übernimmt die Darstellung.

// sampleCommandTree entfernt: Keine Dummy-Daten mehr im Code.
