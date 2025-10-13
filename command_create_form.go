package main

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"sewworkspacemanager/internal/domain"
)

// CommandCreateForm erfasst Daten für einen neuen Befehl und dessen Zielordner.
type CommandCreateForm struct {
	inputs []*textinput.Model
	labels []string
	focus  int
	fieldW int
	help   help.Model
	keys   createKeyMap
	// folder Auswahl (Pfad als Slice Labels)
	folders     [][]string // mögliche Zielpfade (jede ist Liste von Labels der Knoten von Root->Ordner)
	folderIndex int
	// Autovervollständigung für manuellen Ordnerpfad
	folderFlat []string // geflattete Liste vorhandener Ordnerpfade (ohne (Top))
	suggestion string   // aktueller Vorschlag (Suffix, das angehängt würde)
}

const (
	idxProg = iota
	idxArgs
	idxWorkdir
	idxFolderPath
	idxDisplay
	idxDescription
)

type createKeyMap struct {
	Submit     key.Binding
	Close      key.Binding
	Next       key.Binding
	Prev       key.Binding
	NextFolder key.Binding
	PrevFolder key.Binding
	Complete   key.Binding
}

func (k createKeyMap) ShortHelp() []key.Binding { return []key.Binding{k.Submit, k.Close} }
func (k createKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Submit, k.Close, k.Next, k.Prev, k.NextFolder, k.PrevFolder, k.Complete}}
}

func NewCommandCreateForm(root *domain.CommandNode) *CommandCreateForm {
	f := &CommandCreateForm{fieldW: 40}
	f.labels = []string{"Programm", "Argumente", "Arbeitsverzeichnis", "Ordnerpfad (optional)", "Anzeigename", "Beschreibung"}
	f.inputs = make([]*textinput.Model, len(f.labels))
	f.help = help.New()
	f.help.Styles.ShortKey = hotkeyStyle
	f.help.Styles.ShortDesc = textStyle
	f.help.Styles.FullKey = hotkeyStyle
	f.help.Styles.FullDesc = textStyle
	f.keys = createKeyMap{
		Submit:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "Speichern")),
		Close:      key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("Esc", "Abbrechen")),
		Next:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("Tab", "Weiter Feld")),
		Prev:       key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("Shift+Tab", "Zurück Feld")),
		NextFolder: key.NewBinding(key.WithKeys("ctrl+down", "down"), key.WithHelp("↓ / ^+↓", "Ordner+")),
		PrevFolder: key.NewBinding(key.WithKeys("ctrl+up", "up"), key.WithHelp("↑ / ^+↑", "Ordner-")),
		Complete:   key.NewBinding(key.WithKeys("ctrl+right"), key.WithHelp("^+→", "Pfad vervollst.")),
	}
	for i := range f.inputs {
		ti := textinput.New()
		ti.Prompt = ""
		ti.CharLimit = 0
		ti.Width = f.fieldW
		f.inputs[i] = &ti
	}
	f.inputs[idxProg].Placeholder = "git"
	f.inputs[idxArgs].Placeholder = "status"
	f.inputs[idxWorkdir].Placeholder = "C:/ADO/CS2"
	f.inputs[idxFolderPath].Placeholder = "Git/Unterordner"
	f.inputs[idxDisplay].Placeholder = "Neuer Befehl"
	f.inputs[idxDescription].Placeholder = "Kurzbeschreibung"
	f.collectFolders(root)
	f.setFocus(0)
	return f
}

func (f *CommandCreateForm) collectFolders(root *domain.CommandNode) {
	var paths [][]string
	var flat []string
	var walk func(n *domain.CommandNode, path []string, depth int)
	walk = func(n *domain.CommandNode, path []string, depth int) {
		if depth > 2 {
			return
		}
		// Ordner ist jeder Knoten (auch wenn Befehl) mit Children; zusätzlich Root (leeres Label) soll auswählbar gelten => mappe leeres Label auf "(Top-Level)"
		name := n.Label
		if name == "" {
			name = "(Top)"
		}
		curPath := append(path, name)
		if n.Label == "" || len(n.Children) > 0 { // wählbar als Ordner
			// speichere Pfad ohne Root (wenn Root synthetisch)
			store := curPath
			paths = append(paths, store)
		}
		for _, ch := range n.Children {
			walk(ch, curPath, depth+1)
		}
	}
	walk(root, []string{}, 0)
	// Sortiere Pfade alphabetisch nach Joined-Darstellung
	sort.Slice(paths, func(i, j int) bool { return strings.Join(paths[i], "/") < strings.Join(paths[j], "/") })
	sort.Strings(flat)
	f.folders = paths
	f.folderIndex = 0
	f.folderFlat = flat
}

func (f *CommandCreateForm) SetWidth(w int) {
	if w <= 0 {
		return
	}
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

func (f *CommandCreateForm) setFocus(i int) {
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
func (f *CommandCreateForm) next() { f.setFocus((f.focus + 1) % len(f.inputs)) }
func (f *CommandCreateForm) prev() { f.setFocus((f.focus - 1 + len(f.inputs)) % len(f.inputs)) }

func (f *CommandCreateForm) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "tab":
			f.next()
			return nil
		case "shift+tab":
			f.prev()
			return nil
		case "ctrl+down", "down":
			f.folderIndex = (f.folderIndex + 1) % len(f.folders)
			return nil
		case "ctrl+up", "up":
			f.folderIndex = (f.folderIndex - 1 + len(f.folders)) % len(f.folders)
			return nil
		case "ctrl+right":
			// Vorschlag übernehmen
			if f.focus == idxFolderPath && f.suggestion != "" && f.inputs[idxFolderPath] != nil {
				cur := f.inputs[idxFolderPath].Value()
				f.inputs[idxFolderPath].SetValue(cur + f.suggestion)
				f.suggestion = ""
			}
			return nil
		}
	}
	if f.inputs[f.focus] != nil {
		var cmd tea.Cmd
		old := f.inputs[f.focus].Value()
		*f.inputs[f.focus], cmd = f.inputs[f.focus].Update(msg)
		if f.focus == idxFolderPath && f.inputs[idxFolderPath] != nil {
			newVal := f.inputs[idxFolderPath].Value()
			if newVal != old {
				f.computeSuggestion(newVal)
			}
		}
		return cmd
	}
	return nil
}

// computeSuggestion berechnet anhand vorhandener Ordnerpfade einen Autocomplete-Vorschlag.
func (f *CommandCreateForm) computeSuggestion(cur string) {
	f.suggestion = ""
	if strings.TrimSpace(cur) == "" {
		return
	}
	// Versuche längsten gemeinsamen Kandidaten unter folderFlat zu finden
	candidates := make([]string, 0, len(f.folderFlat))
	for _, fp := range f.folderFlat {
		if strings.HasPrefix(strings.ToLower(fp), strings.ToLower(cur)) {
			candidates = append(candidates, fp)
		}
	}
	if len(candidates) == 0 {
		return
	}
	// Finde gemeinsamen Präfix (case-insensitive) der Kandidaten
	prefix := candidates[0]
	for _, c := range candidates[1:] {
		prefix = commonPrefixCaseFold(prefix, c)
		if prefix == "" {
			break
		}
	}
	if len(prefix) <= len(cur) {
		return
	}
	// Begrenze Tiefe auf 2 Ordner Ebenen
	segs := strings.Split(prefix, "/")
	if len(segs) > 2 { // wir erlauben max 2 segs (entspricht zwei Ordner Ebenen unter Root)
		segs = segs[:2]
		prefix = strings.Join(segs, "/")
	}
	if len(prefix) > len(cur) {
		f.suggestion = prefix[len(cur):]
	}
}

// commonPrefixCaseFold liefert gemeinsamen Präfix (case-insensitive Vergleich, gibt ursprüngliche Zeichen des ersten Strings zurück)
func commonPrefixCaseFold(a, b string) string {
	la, lb := len(a), len(b)
	n := la
	if lb < n {
		n = lb
	}
	i := 0
	for i < n {
		if strings.EqualFold(a[i:i+1], b[i:i+1]) {
			i++
		} else {
			break
		}
	}
	return a[:i]
}

// Values liefert die eingegebenen Daten.
func (f *CommandCreateForm) Values() (prog string, args []string, workdir string, display string, desc string, folderPath []string) {
	if f.inputs[idxProg] != nil {
		prog = strings.TrimSpace(f.inputs[idxProg].Value())
		if prog == "" {
			prog = f.inputs[idxProg].Placeholder
		}
	}
	if f.inputs[idxArgs] != nil {
		raw := f.inputs[idxArgs].Value()
		if strings.TrimSpace(raw) == "" {
			raw = f.inputs[idxArgs].Placeholder
		}
		args = parseArgs(raw)
	}
	if f.inputs[idxWorkdir] != nil {
		workdir = strings.TrimSpace(f.inputs[idxWorkdir].Value())
		if workdir == "" {
			workdir = f.inputs[idxWorkdir].Placeholder
		}
	}
	if f.inputs[idxDisplay] != nil {
		display = strings.TrimSpace(f.inputs[idxDisplay].Value())
		if display == "" {
			display = f.inputs[idxDisplay].Placeholder
		}
	}
	if f.inputs[idxDescription] != nil {
		desc = strings.TrimSpace(f.inputs[idxDescription].Value())
		if desc == "" {
			desc = f.inputs[idxDescription].Placeholder
		}
	}
	// Ordnerpfad: wenn gesetzt, überschreibe Auswahl
	manual := ""
	if f.inputs[idxFolderPath] != nil {
		manual = strings.TrimSpace(f.inputs[idxFolderPath].Value())
	}
	if manual != "" {
		// Split auf '/'
		parts := strings.Split(manual, "/")
		cleaned := make([]string, 0, len(parts)+1)
		cleaned = append(cleaned, "(Top)")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			cleaned = append(cleaned, p)
			if len(cleaned)-1 >= 3 {
				break
			} // max zwei Ordner-Ebenen (gesamt Tiefe 3 mit Root)
		}
		folderPath = cleaned
	} else if len(f.folders) > 0 {
		folderPath = f.folders[f.folderIndex]
	}
	return
}

func (f *CommandCreateForm) View() string {
	lines := []string{RenderDialogHeader("Neu", f.fieldW+18), ""}
	for i, ti := range f.inputs {
		label := f.labels[i]
		valView := ti.View()
		if i == idxFolderPath && f.suggestion != "" {
			// Vorschlag farblich andeuten (hell, fallback falls keine definierte Graufarbe existiert)
			valView += textStyle.Copy().Faint(true).Render(f.suggestion)
		}
		lines = append(lines, label+": "+valView)
	}
	// Folder Auswahl
	curFolder := "(keine)"
	if len(f.folders) > 0 {
		curFolder = strings.Join(f.folders[f.folderIndex], "/")
	}
	lines = append(lines, "")
	manual := ""
	if f.inputs[idxFolderPath] != nil {
		manual = strings.TrimSpace(f.inputs[idxFolderPath].Value())
	}
	lines = append(lines, "Zielordner (Auswahl): "+curFolder)
	lines = append(lines, "(^+↑/^+↓ Ordner wechseln oder Feld 'Ordnerpfad' verwenden)")
	if manual != "" {
		lines = append(lines, "Manueller Pfad aktiv: "+manual)
	}
	lines = append(lines, "")
	lines = append(lines, strings.Split(f.help.View(f.keys), "\n")...)
	return strings.Join(lines, "\n")
}
