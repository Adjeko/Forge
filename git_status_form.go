package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// GitStatusForm erfasst nur einen Repository-Pfad für 'git status'.
type GitStatusForm struct {
	input  *textinput.Model
	fieldW int
	help   help.Model
	keys   gitStatusKeyMap
}

type gitStatusKeyMap struct {
	Submit key.Binding
	Close  key.Binding
}

func (k gitStatusKeyMap) ShortHelp() []key.Binding  { return []key.Binding{k.Submit, k.Close} }
func (k gitStatusKeyMap) FullHelp() [][]key.Binding { return [][]key.Binding{{k.Submit, k.Close}} }

func NewGitStatusForm() *GitStatusForm {
	f := &GitStatusForm{fieldW: 48}
	f.help = help.New()
	f.help.Styles.ShortKey = hotkeyStyle
	f.help.Styles.ShortDesc = textStyle
	f.help.Styles.FullKey = hotkeyStyle
	f.help.Styles.FullDesc = textStyle
	f.keys = gitStatusKeyMap{
		Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "Ausführen")),
		Close:  key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("Esc", "Schließen")),
	}
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = "C:/ADO/CS2"
	ti.CharLimit = 0
	ti.Width = f.fieldW
	ti.Focus()
	f.input = &ti
	return f
}

func (f *GitStatusForm) SetWidth(w int) {
	if w <= 0 {
		return
	}
	field := w - 18
	if field < 20 {
		field = 20
	}
	f.fieldW = field
	if f.input != nil {
		f.input.Width = field
	}
	f.help.Width = w
}

func (f *GitStatusForm) Update(msg tea.Msg) tea.Cmd {
	if f.input == nil {
		return nil
	}
	var cmd tea.Cmd
	*f.input, cmd = f.input.Update(msg)
	return cmd
}

// Value liefert den bereinigten Pfad.
func (f *GitStatusForm) Value() string {
	if f.input == nil {
		return ""
	}
	v := strings.TrimSpace(f.input.Value())
	if v == "" {
		v = strings.TrimSpace(f.input.Placeholder)
	}
	return v
}

// Validate prüft, ob Pfad existiert und ein .git Verzeichnis enthält.
func (f *GitStatusForm) Validate() (string, error) {
	path := f.Value()
	if path == "" {
		return "", ErrInvalidRepoPath("Pfad leer")
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return "", ErrInvalidRepoPath("Pfad existiert nicht oder ist kein Ordner")
	}
	gitDir := filepath.Join(path, ".git")
	gi, err := os.Stat(gitDir)
	if err != nil || !gi.IsDir() {
		return "", ErrInvalidRepoPath("Kein .git Ordner gefunden")
	}
	return path, nil
}

func (f *GitStatusForm) View() string {
	lines := []string{RenderDialogHeader("Git Status", f.fieldW+18), ""}
	lines = append(lines, "Repository Pfad: "+f.input.View())
	lines = append(lines, "")
	lines = append(lines, strings.Split(f.help.View(f.keys), "\n")...)
	return strings.Join(lines, "\n")
}

// Fehler-Typ für Validierung
type invalidRepoPathError struct{ msg string }

func (e invalidRepoPathError) Error() string { return e.msg }
func ErrInvalidRepoPath(msg string) error    { return invalidRepoPathError{msg: msg} }
