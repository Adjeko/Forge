package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CommandForm ist eine kleine Bubble-Komponente, die einen Befehlsnamen, seine Argumente und ein Arbeitsverzeichnis erfasst.
type CommandForm struct {
	inputs []*textinput.Model
	labels []string
	focus  int
	fieldW int
	help   help.Model
	keys   formKeyMap
}

// formKeyMap definiert formular-spezifische Tastenkürzel für Hilfe-Darstellung.
type formKeyMap struct {
	Submit key.Binding
	Close  key.Binding
	Next   key.Binding
	Prev   key.Binding
}

func (k formKeyMap) ShortHelp() []key.Binding { return []key.Binding{k.Submit, k.Close} }
func (k formKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Submit, k.Close},
		{k.Next, k.Prev},
	}
}

func NewCommandForm() *CommandForm {
	f := &CommandForm{
		inputs: make([]*textinput.Model, 3),
		labels: []string{"Befehl", "Argumente", "Arbeitsverzeichnis"},
		focus:  0,
		fieldW: 40,
	}
	// Hilfe + Tasten
	f.help = help.New()
	// Hilfe-Styling an Haupt-Footer-Stile angleichen (falls verfügbar)
	f.help.Styles.ShortKey = hotkeyStyle
	f.help.Styles.ShortDesc = textStyle
	f.help.Styles.FullKey = hotkeyStyle
	f.help.Styles.FullDesc = textStyle
	f.keys = formKeyMap{
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "Ausführen"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("Esc", "Schließen"),
		),
		Next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "Nächstes Feld"),
		),
		Prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("Shift+Tab", "Vorheriges Feld"),
		),
	}
	// Eingabefelder erzeugen
	for i := 0; i < len(f.inputs); i++ {
		ti := textinput.New()
		ti.Prompt = ""
		ti.CharLimit = 0 // unbegrenzt
		ti.Placeholder = ""
		width := f.fieldW
		if width < 10 {
			width = 10
		}
		ti.Width = width
		f.inputs[i] = &ti
	}
	// Defaults: häufig genutzter Fall
	f.inputs[0].Placeholder = "git"
	f.inputs[1].Placeholder = "status"
	f.inputs[2].Placeholder = "C:/ADO/CS2"
	f.setFocus(0)
	return f
}

// SetWidth passt die Eingabe-Breiten an den verfügbaren Inhaltsbereich an.
func (f *CommandForm) SetWidth(w int) {
	if w <= 0 {
		return
	}
	// Platz für Label und Doppelpunkt lassen; Mindestbreite sicherstellen
	field := w - 18
	if field < 10 {
		field = 10
	}
	f.fieldW = field
	f.help.Width = w
	for _, ti := range f.inputs {
		if ti != nil {
			ti.Width = field
		}
	}
}

func (f *CommandForm) setFocus(i int) {
	if i < 0 || i >= len(f.inputs) {
		return
	}
	for idx, ti := range f.inputs {
		if ti == nil {
			continue
		}
		if idx == i {
			ti.Focus()
		} else {
			ti.Blur()
		}
	}
	f.focus = i
}

func (f *CommandForm) next() { f.setFocus((f.focus + 1) % len(f.inputs)) }
func (f *CommandForm) prev() { f.setFocus((f.focus - 1 + len(f.inputs)) % len(f.inputs)) }

func (f *CommandForm) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "tab":
			f.next()
			return nil
		case "shift+tab":
			f.prev()
			return nil
		}
	}
	// Nur fokussiertes Feld aktualisieren
	if f.inputs[f.focus] != nil {
		var cmd tea.Cmd
		*f.inputs[f.focus], cmd = f.inputs[f.focus].Update(msg)
		return cmd
	}
	return nil
}

func (f *CommandForm) View() string {
	lines := make([]string, 0, 8)
	// Header Zeile
	lines = append(lines, RenderDialogHeader("Befehl", f.fieldW+18))
	lines = append(lines, "")
	for i, ti := range f.inputs {
		label := f.labels[i]
		line := label + ": " + ti.View()
		lines = append(lines, line)
	}
	// Hilfe ähnlich wie Haupt-Footer gestylt (hotkeyStyle/textStyle), falls verfügbar
	// Hilfe-View bauen: wir können hier nicht auf Breite vertrauen; Container setzt Breite
	helpView := f.help.View(f.keys)
	lines = append(lines, "")
	// Mehrzeilige Hilfe zeilenweise anhängen
	lines = append(lines, strings.Split(helpView, "\n")...)
	return strings.Join(lines, "\n")
}

// Values liefert die Formularwerte. Argumente unterstützen einfache Maskierung von Leerzeichen mit Quotes.
func (f *CommandForm) Values() (cmd string, args []string, dir string) {
	if f.inputs[0] != nil {
		cmd = strings.TrimSpace(f.inputs[0].Value())
		if cmd == "" {
			cmd = strings.TrimSpace(f.inputs[0].Placeholder)
		}
	}
	if f.inputs[1] != nil {
		raw := f.inputs[1].Value()
		if strings.TrimSpace(raw) == "" {
			raw = f.inputs[1].Placeholder
		}
		args = parseArgs(raw)
	}
	if f.inputs[2] != nil {
		dir = strings.TrimSpace(f.inputs[2].Value())
		if dir == "" {
			dir = strings.TrimSpace(f.inputs[2].Placeholder)
		}
	}
	return
}

// parseArgs teilt String in Argumente mit einfacher Quote-Behandlung (einfach und doppelt).
func parseArgs(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var (
		args   []string
		b      strings.Builder
		inSQ   bool
		inDQ   bool
		escape bool
	)
	flush := func() {
		if b.Len() > 0 {
			args = append(args, b.String())
			b.Reset()
		}
	}
	for _, r := range s {
		switch {
		case escape:
			b.WriteRune(r)
			escape = false
		case r == '\\':
			escape = true
		case r == '\'' && !inDQ:
			inSQ = !inSQ
		case r == '"' && !inSQ:
			inDQ = !inDQ
		case r == ' ' && !inSQ && !inDQ:
			flush()
		default:
			b.WriteRune(r)
		}
	}
	flush()
	return args
}
