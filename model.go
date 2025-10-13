package main

import (
	"bufio"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sewworkspacemanager/internal/domain"
	"sewworkspacemanager/internal/persist"
)

// Gemeinsame Styles
var (
	// Footer: Hotkeys in hellgrau, Beschreibungen dunkler (Hotkeys fett)
	hotkeyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Bold(true)
	textStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	contentStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	footerLineStyle = lipgloss.NewStyle()
	// Modal Rahmen-Stil fixiert auf helleres SEW Rot
	modalBorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4A5A"))
	// Fehlerstil (stderr) in Rot
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E30018"))
	// Erfolgsstil in Grün
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#33BB33"))
)

// PullRequest enthält Minimaldaten für einen PR-Eintrag.
type PullRequest struct {
	Number string
	Title  string
	Author string
	Branch string
	Status string // z.B. OPEN, MERGED, DRAFT
}

// (veraltetes Duplikat entfernt unten)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// (veraltetes Duplikat entfernt)

// ensureBefehleRoot lädt die Befehle-Wurzel falls sie fehlt.
func (m *Model) ensureBefehleRoot() {
	if m.rootBefehle == nil {
		br, _, err := persist.LoadDualCommandTrees("commands.toml")
		if err != nil || br == nil {
			br = &domain.CommandNode{Label: "", Expanded: true}
		}
		m.rootBefehle = br
	}
}

// ensureAblaeufeRoot lädt die Abläufe-Wurzel falls sie fehlt.
func (m *Model) ensureAblaeufeRoot() {
	if m.rootAblaeufe == nil {
		_, ar, err := persist.LoadDualCommandTrees("commands.toml")
		if err != nil || ar == nil {
			ar = &domain.CommandNode{Label: "", Expanded: true}
		}
		m.rootAblaeufe = ar
	}
}

// (ensureDualRoots entfernt – granularer Helper vorhanden)

// ensureStatusRoots lädt die Status-Wurzel und ergänzt Befehle/Abläufe falls fehlend.
func (m *Model) ensureStatusRoots() {
	if m.rootStatus == nil || m.rootBefehle == nil || m.rootAblaeufe == nil {
		br, ar, sr, err := persist.LoadAllCommandGroups("commands.toml")
		if m.rootBefehle == nil {
			if err != nil || br == nil {
				br = &domain.CommandNode{Label: "", Expanded: true}
			}
			m.rootBefehle = br
		}
		if m.rootAblaeufe == nil {
			if err != nil || ar == nil {
				ar = &domain.CommandNode{Label: "", Expanded: true}
			}
			m.rootAblaeufe = ar
		}
		if m.rootStatus == nil {
			if err != nil || sr == nil {
				sr = &domain.CommandNode{Label: "", Expanded: true}
			}
			m.rootStatus = sr
		}
	}
}

type Model struct {
	width            int
	height           int
	showHelp         bool
	showModal        bool
	form             *CommandForm
	showTree         bool
	showFlows        bool
	showStatus       bool
	showCreate       bool
	showFlowsCreate  bool
	showStatusCreate bool
	createForm       *CommandCreateForm
	rootBefehle      *domain.CommandNode
	rootAblaeufe     *domain.CommandNode
	rootStatus       *domain.CommandNode
	flowsCreateForm  *CommandCreateForm
	statusCreateForm *CommandCreateForm
	// Ausgabebereich (Suchleiste entfernt)
	termOutput []string
	// Status-Spalten-Icons (einfache Platzhalterliste)
	statusIcons []string
	// Fortschritts-Kacheln unten
	tiles       []ProgressTile
	currentPage string
	// Git Status Dialog
	showGitStatus bool
	gitStatusForm *GitStatusForm
	// Scrollbarer Viewport für Terminalausgabe
	vp viewport.Model
	// Streaming-Prozessstatus
	streaming  bool
	scannerOut *bufio.Scanner
	scannerErr *bufio.Scanner
	tree       *CommandTree
	doneOut    bool
	doneErr    bool
	cmd        *exec.Cmd
	// Hilfe und Key Bindings
	help help.Model
	keys keyMap
	// Pull-Requests Daten
	pullRequests []PullRequest
	loadingPRs   bool
	// Streaming eines einzelnen Befehls (Multi-Task Fortschritt entfernt)
}

// keyMap definiert Tastenkombinationen und implementiert help.KeyMap.
type keyMap struct {
	ToggleHelp       key.Binding
	PageDashboard    key.Binding
	PagePullRequests key.Binding
	Quit             key.Binding
	Esc              key.Binding
	Enter            key.Binding
	Tree             key.Binding
	Flows            key.Binding
	Status           key.Binding
	Create           key.Binding
	AddTile          key.Binding
	Up               key.Binding
	Down             key.Binding
	PageUp           key.Binding
	PageDown         key.Binding
	GitStatus        key.Binding
	FetchPRs         key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	// Entfernt Quit aus der kompakten Hilfe; Quit erscheint nur in der ausgeklappten Hilfe
	return []key.Binding{k.ToggleHelp, k.PageDashboard, k.PagePullRequests, k.FetchPRs, k.Tree, k.Flows, k.Status}
}
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleHelp, k.PageDashboard, k.PagePullRequests, k.FetchPRs, k.Quit, k.Tree, k.Flows, k.Status, k.Esc, k.Enter, k.AddTile},
		{k.Up, k.Down, k.PageUp, k.PageDown},
	}
}

// NewModel erzeugt das initiale Model.
func NewModel() Model {
	m := Model{currentPage: "Dashboard"}
	// Initial einen Viewport mit Größe 0 ohne Rahmen erzeugen
	m.vp = viewport.New(0, 0)
	m.vp.Style = lipgloss.NewStyle() // kein Frame/Rahmen
	// Anfangs-Status-Icons (später dynamisch möglich)
	m.statusIcons = []string{"✅", "⚠", "⌛"}
	// Platzhalter-Fortschritts-Kacheln
	m.tiles = []ProgressTile{
		NewProgressTile("1", "Build"),
		NewProgressTile("2", "Tests"),
	}
	// Hilfe und Keys
	m.help = help.New()
	m.keys = keyMap{
		ToggleHelp: key.NewBinding(
			key.WithKeys("ctrl+y"),
			// Beschriftung verkürzt auf nur "Hilfe"
			key.WithHelp("^+Y", "Hilfe"),
		),
		PageDashboard: key.NewBinding(
			key.WithKeys("ctrl+1"),
			key.WithHelp("^+1", "Dashboard"),
		),
		PagePullRequests: key.NewBinding(
			key.WithKeys("ctrl+2"),
			key.WithHelp("^+2", "Pull Requests"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("^+C", "Beenden"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "Schließen"),
		),
		Tree: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("^+B", "Befehle"),
		),
		Flows: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("^+A", "Abläufe"),
		),
		Status: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("^+S", "Status"),
		),
		Create: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("^+N", "Befehl neu"),
		),
		AddTile: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("^+T", "Tile hinzufügen"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "Git Status Ausführen"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "Scrollen"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "Scrollen"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("PgUp", "Hoch"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("PgDn", "Runter"),
		),
		GitStatus: key.NewBinding(
			key.WithKeys("ctrl+g"),
			key.WithHelp("^+G", "Git Status"),
		),
		FetchPRs: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("^+P", "PRs laden"),
		),
	}
	return m
}

// Init erfüllt tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(time.Time) tea.Msg { return progressTickMsg{} })
}

// Update verarbeitet Tastatur- und Fenster-Events.

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = recalcViewport(m)
		return m, nil
	case streamLineMsg:
		// Gestreamte Zeile anhängen, Viewport aktualisieren und nächste Zeile anfordern
		var styled string
		if msg.isErr {
			styled = errorStyle.Render("[stderr] ") + errorStyle.Render(msg.line)
		} else {
			styled = contentStyle.Render(msg.line)
		}
		// Einzelner Streaming-Puffer
		m.termOutput = append(m.termOutput, styled)
		m.vp.SetContent(strings.Join(m.termOutput, "\n"))
		m.vp.GotoBottom()
		return m, m.readNextLine(msg.isErr)
	case streamDoneMsg:
		// Streaming-Status bereinigen und ggf. Fehler anhängen
		if msg.isErr {
			m.doneErr = true
			m.scannerErr = nil
		} else {
			m.doneOut = true
			m.scannerOut = nil
		}
		if m.doneOut && m.doneErr {
			m.streaming = false
			m.cmd = nil
		}
		if msg.err != "" {
			m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+msg.err))
			m.vp.SetContent(strings.Join(m.termOutput, "\n"))
			m.vp.GotoBottom()
		}
		return m, nil
	case tea.KeyMsg:
		// Wenn ein Modal offen ist, Eingabe zuerst an das Formular leiten
		if m.showModal {
			// Schließen und Absenden innerhalb des Modal-Kontextes behandeln
			switch {
			case key.Matches(msg, m.keys.Esc):
				m.showModal = false
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.form != nil && !m.streaming {
					name, args, dir := m.form.Values()
					if name != "" {
						startCmd, err := StartStreamingCommand(&m, name, args, dir)
						m.showModal = false
						if err != nil {
							m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+err.Error()))
							m.vp.SetContent(strings.Join(m.termOutput, "\n"))
							m.vp.GotoBottom()
							return m, nil
						}
						return m, startCmd
					}
				}
				return m, nil
			}
			// Andernfalls Formular aktualisieren und keine Tastendrücke an Viewport durchreichen
			if m.form != nil {
				cmd := m.form.Update(msg)
				return m, cmd
			}
			return m, nil
		} else if m.showCreate {
			// Erstellen-Dialog
			switch {
			case key.Matches(msg, m.keys.Esc) || key.Matches(msg, m.keys.Create):
				m.showCreate = false
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.createForm != nil {
					prog, args, workdir, display, desc, folderPath := m.createForm.Values()
					if prog != "" && display != "" {
						// Lade Befehle-Wurzel falls fehlend
						m.ensureBefehleRoot()
						domain.AddNodeToTree(m.rootBefehle, folderPath, &domain.CommandNode{Label: display, Description: desc, Command: prog, Args: args, WorkDir: workdir})
						if err := persist.SaveDualCommandTrees(m.rootBefehle, m.rootAblaeufe, "commands.toml"); err != nil {
							m.termOutput = append(m.termOutput, errorStyle.Render("[error] Speichern fehlgeschlagen: "+err.Error()))
						} else {
							m.termOutput = append(m.termOutput, successStyle.Render("["+display+"] wurde erfolgreich angelegt"))
						}
						m.vp.SetContent(strings.Join(m.termOutput, "\n"))
						m.vp.GotoBottom()
						m.showCreate = false
					}
					return m, nil
				}
			}
			if m.createForm != nil {
				cmd := m.createForm.Update(msg)
				return m, cmd
			}
			return m, nil
		} else if m.showFlowsCreate {
			// Abläufe Erstellen-Dialog (eigener Fokus)
			switch {
			case key.Matches(msg, m.keys.Esc):
				m.showFlowsCreate = false
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.flowsCreateForm != nil {
					prog, args, workdir, display, desc, folderPath := m.flowsCreateForm.Values()
					if prog != "" && display != "" {
						m.ensureAblaeufeRoot()
						domain.AddNodeToTree(m.rootAblaeufe, folderPath, &domain.CommandNode{Label: display, Description: desc, Command: prog, Args: args, WorkDir: workdir})
						// Speichern
						if err := persist.SaveDualCommandTrees(m.rootBefehle, m.rootAblaeufe, "commands.toml"); err != nil {
							m.termOutput = append(m.termOutput, errorStyle.Render("[error] Speichern fehlgeschlagen: "+err.Error()))
						} else {
							m.termOutput = append(m.termOutput, successStyle.Render("[Ablauf "+display+"] wurde erfolgreich angelegt"))
						}
						m.vp.SetContent(strings.Join(m.termOutput, "\n"))
						m.vp.GotoBottom()
						m.showFlowsCreate = false
					}
					return m, nil
				}
			}
			if m.flowsCreateForm != nil {
				cmd := m.flowsCreateForm.Update(msg)
				return m, cmd
			}
			return m, nil
		} else if m.showStatusCreate {
			// Status Erstellen-Dialog
			switch {
			case key.Matches(msg, m.keys.Esc):
				m.showStatusCreate = false
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.statusCreateForm != nil {
					prog, args, workdir, display, desc, folderPath := m.statusCreateForm.Values()
					if prog != "" && display != "" {
						m.ensureStatusRoots()
						domain.AddNodeToTree(m.rootStatus, folderPath, &domain.CommandNode{Label: display, Description: desc, Command: prog, Args: args, WorkDir: workdir})
						if err := persist.SaveAllCommandGroups(m.rootBefehle, m.rootAblaeufe, m.rootStatus, "commands.toml"); err != nil {
							m.termOutput = append(m.termOutput, errorStyle.Render("[error] Speichern fehlgeschlagen: "+err.Error()))
						} else {
							m.termOutput = append(m.termOutput, successStyle.Render("[Status "+display+"] wurde erfolgreich angelegt"))
						}
						m.vp.SetContent(strings.Join(m.termOutput, "\n"))
						m.vp.GotoBottom()
						m.showStatusCreate = false
					}
					return m, nil
				}
			}
			if m.statusCreateForm != nil {
				cmd := m.statusCreateForm.Update(msg)
				return m, cmd
			}
			return m, nil
		} else if m.showTree {
			// Befehle-Baum Dialog Eingabeverarbeitung
			switch {
			case key.Matches(msg, m.keys.Esc) || key.Matches(msg, m.keys.Tree):
				m.showTree = false
				return m, nil
			case key.Matches(msg, m.keys.Create):
				m.showCreate = true
				m.ensureBefehleRoot()
				m.createForm = NewCommandCreateForm(m.rootBefehle)
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.tree != nil {
					if cmdNode := m.tree.SelectedCommand(); cmdNode != nil && !m.streaming {
						if cmdNode.Command == "git" && len(cmdNode.Args) == 1 && cmdNode.Args[0] == "status" {
							m.showGitStatus = true
							m.gitStatusForm = NewGitStatusForm()
							m.showTree = false
							return m, nil
						}
						startCmd, err := StartStreamingCommand(&m, cmdNode.Command, cmdNode.Args, cmdNode.WorkDir)
						m.showTree = false
						if err != nil {
							m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+err.Error()))
							m.vp.SetContent(strings.Join(m.termOutput, "\n"))
							m.vp.GotoBottom()
							return m, nil
						}
						return m, startCmd
					}
				}
			}
			if m.tree != nil {
				var tcmd tea.Cmd
				m.tree, tcmd = m.tree.Update(msg)
				return m, tcmd
			}
			return m, nil
		} else if m.showFlows {
			// Abläufe Baum Dialog
			switch {
			case key.Matches(msg, m.keys.Esc) || key.Matches(msg, m.keys.Flows):
				m.showFlows = false
				return m, nil
			case key.Matches(msg, m.keys.Create):
				m.showFlowsCreate = true
				if m.rootAblaeufe == nil {
					_, ar, err := persist.LoadDualCommandTrees("commands.toml")
					if err != nil || ar == nil {
						ar = &domain.CommandNode{Label: "", Expanded: true}
					}
					m.rootAblaeufe = ar
				}
				m.flowsCreateForm = NewCommandCreateForm(m.rootAblaeufe)
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.tree != nil {
					if cmdNode := m.tree.SelectedCommand(); cmdNode != nil && !m.streaming {
						startCmd, err := StartStreamingCommand(&m, cmdNode.Command, cmdNode.Args, cmdNode.WorkDir)
						m.showFlows = false
						if err != nil {
							m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+err.Error()))
							m.vp.SetContent(strings.Join(m.termOutput, "\n"))
							m.vp.GotoBottom()
							return m, nil
						}
						return m, startCmd
					}
				}
			}
			if m.tree != nil {
				var tcmd tea.Cmd
				m.tree, tcmd = m.tree.Update(msg)
				return m, tcmd
			}
			return m, nil
		} else if m.showStatus {
			// Status Baum Dialog Verarbeitung (analog zu Befehle/Abläufe)
			switch {
			case key.Matches(msg, m.keys.Esc) || key.Matches(msg, m.keys.Status):
				m.showStatus = false
				return m, nil
			case key.Matches(msg, m.keys.Create):
				m.showStatusCreate = true
				if m.rootStatus == nil {
					br, ar, sr, err := persist.LoadAllCommandGroups("commands.toml")
					if err != nil {
						br = &domain.CommandNode{Label: "", Expanded: true}
						ar = &domain.CommandNode{Label: "", Expanded: true}
						sr = &domain.CommandNode{Label: "", Expanded: true}
					}
					if br == nil {
						br = &domain.CommandNode{Label: "", Expanded: true}
					}
					if ar == nil {
						ar = &domain.CommandNode{Label: "", Expanded: true}
					}
					if sr == nil {
						sr = &domain.CommandNode{Label: "", Expanded: true}
					}
					m.rootBefehle, m.rootAblaeufe, m.rootStatus = br, ar, sr
				}
				m.statusCreateForm = NewCommandCreateForm(m.rootStatus)
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.tree != nil {
					if cmdNode := m.tree.SelectedCommand(); cmdNode != nil && !m.streaming {
						startCmd, err := StartStreamingCommand(&m, cmdNode.Command, cmdNode.Args, cmdNode.WorkDir)
						m.showStatus = false
						if err != nil {
							m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+err.Error()))
							m.vp.SetContent(strings.Join(m.termOutput, "\n"))
							m.vp.GotoBottom()
							return m, nil
						}
						return m, startCmd
					}
				}
			}
			if m.tree != nil {
				var tcmd tea.Cmd
				m.tree, tcmd = m.tree.Update(msg)
				return m, tcmd
			}
			return m, nil
		} else if m.showGitStatus {
			switch {
			case key.Matches(msg, m.keys.Esc) || key.Matches(msg, m.keys.GitStatus):
				m.showGitStatus = false
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				if m.gitStatusForm != nil && !m.streaming {
					repo, err := m.gitStatusForm.Validate()
					if err != nil {
						m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+err.Error()))
						m.vp.SetContent(strings.Join(m.termOutput, "\n"))
						m.vp.GotoBottom()
						return m, nil
					}
					m.termOutput = nil
					m.vp.SetContent("")
					startCmd, serr := StartStreamingCommand(&m, "git", []string{"status"}, repo)
					m.showGitStatus = false
					if serr != nil {
						m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+serr.Error()))
						m.vp.SetContent(strings.Join(m.termOutput, "\n"))
						m.vp.GotoBottom()
						return m, nil
					}
					return m, startCmd
				}
			}
			if m.gitStatusForm != nil {
				cmd := m.gitStatusForm.Update(msg)
				return m, cmd
			}
			return m, nil
		}
		// Hotkeys zuerst behandeln, damit sie nicht im Eingabefeld landen
		var cmds []tea.Cmd
		switch {
		case key.Matches(msg, m.keys.PageDashboard) || msg.String() == "1":
			m.currentPage = "Dashboard"
			m = recalcViewport(m)
			return m, nil
		case key.Matches(msg, m.keys.PagePullRequests) || msg.String() == "2":
			m.currentPage = "Pull Requests"
			m = recalcViewport(m)
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			if m.streaming {
				// Während Streaming nicht sofort beenden (optional: erlauben?) – hier belassen.
			}

			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleHelp):
			m.help.ShowAll = !m.help.ShowAll
			// Layout ändert sich weil sich Footer-Höhe ändert
			m = recalcViewport(m)
			// Weiterleiten zum Viewport falls dort verarbeitet
		case key.Matches(msg, m.keys.Esc):
			if m.showModal {
				m.showModal = false
			} else if m.help.ShowAll {
				m.help.ShowAll = false
			}
			// weiter
		case key.Matches(msg, m.keys.Create):
			// globales Create deaktiviert (nur im Befehle/Abläufe Dialog über Tree-KeyMap genutzt)
			break
		case key.Matches(msg, m.keys.Tree):
			m.showTree = !m.showTree
			if m.showTree {
				m.ensureBefehleRoot()
				m.tree = NewCommandTreeWithTitle(m.rootBefehle, "Befehle")
				if m.tree != nil {
					m.tree.keys.Create = key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("^+N", "Befehl neu"))
				}
			}
			return m, nil
		case key.Matches(msg, m.keys.Flows):
			m.showFlows = !m.showFlows
			if m.showFlows {
				m.ensureAblaeufeRoot()
				m.tree = NewCommandTreeWithTitle(m.rootAblaeufe, "Abläufe")
				if m.tree != nil {
					m.tree.keys.Create = key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("^+N", "Ablauf neu"))
				}
			}
			return m, nil
		case key.Matches(msg, m.keys.Status):
			m.showStatus = !m.showStatus
			if m.showStatus {
				m.ensureStatusRoots()
				m.tree = NewCommandTreeWithTitle(m.rootStatus, "Status")
				if m.tree != nil {
					m.tree.keys.Create = key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("^+N", "Status neu"))
				}
			}
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			// Standardaktion bleibt: git status im Default-Pfad
			if !m.streaming {
				m.termOutput = nil
				m.vp.SetContent("")
				startCmd, err := StartStreamingCommand(&m, "git", []string{"status"}, "C:/ADO/CS2")
				if err != nil {
					m.termOutput = append(m.termOutput, errorStyle.Render("[error] "+err.Error()))
					m.vp.SetContent(strings.Join(m.termOutput, "\n"))
					m.vp.GotoBottom()
					break
				}
				cmds = append(cmds, startCmd)
			}
		case key.Matches(msg, m.keys.FetchPRs):
			// Platzhalter-Ladevorgang: simulieren durch Leeren und Dummy-Einträge
			if !m.loadingPRs {
				m.loadingPRs = true
				m.pullRequests = []PullRequest{
					{Number: "#12", Title: "Feature: Login Flow", Author: "alice", Branch: "feature/login", Status: "OPEN"},
					{Number: "#15", Title: "Fix: Null Pointer", Author: "bob", Branch: "bugfix/npe", Status: "DRAFT"},
					{Number: "#18", Title: "Chore: Update deps", Author: "carol", Branch: "chore/deps", Status: "OPEN"},
				}
				m.loadingPRs = false
				// Noch kein asynchroner Befehl; später gh CLI integrierbar
				m = recalcViewport(m)
			}
			return m, nil
		case key.Matches(msg, m.keys.GitStatus):
			if !m.streaming && !m.showGitStatus { // Dialog öffnen
				m.showGitStatus = true
				m.gitStatusForm = NewGitStatusForm()
			}
		case key.Matches(msg, m.keys.AddTile):
			if len(m.tiles) < 10 {
				id := domain.Itoa(len(m.tiles) + 1)
				label := "Task" + id
				m.tiles = append(m.tiles, NewProgressTile(id, label))
			}
		case msg.String() == "backspace":
			// ignoriert (kein Eingabefeld mehr)
		case strings.HasPrefix(msg.String(), "ctrl+"):
			// früherer Task-Umschaltung Platzhalter – keine Aktion
		}
		// Keine freie Texteingabe mehr (Suchleiste entfernt)
		// Taste immer an Viewport weitergeben um Scrollen zu ermöglichen
		var vpcmd tea.Cmd
		m.vp, vpcmd = m.vp.Update(msg)
		if vpcmd != nil {
			cmds = append(cmds, vpcmd)
		}
		if len(cmds) > 0 {
			return m, tea.Batch(cmds...)
		}
		return m, nil
	case tea.MouseMsg:
		// Maus-Events zum Scrollen weiterleiten
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd
	case progressTickMsg:
		// Jede Kachel-Prozentzahl aktualisieren und nächsten Tick einplanen
		var cmd tea.Cmd
		for i := range m.tiles {
			// Tick wiederverwenden um nächste Planung zu bekommen (nur einen Cmd sammeln)
			c := m.tiles[i].Tick()
			if cmd == nil {
				cmd = c
			}
		}
		// Viewport wird nicht neu berechnet, Anzeige zeigt aktualisierte Balken
		return m, cmd
	}
	return m, nil
}

// padToLines stellt sicher, dass der gerenderte Block genau n Zeilen hat (auffüllen oder kürzen).
func padToLines(s string, n int) string {
	if n <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		return strings.Join(lines[:n], "\n")
	}
	if len(lines) < n {
		return strings.Join(append(lines, make([]string, n-len(lines))...), "\n")
	}
	return s
}

// insertAt legt segment an Position x in base und hält Gesamtbreite w. Annahme: base-Länge == w.
func insertAt(base, segment string, x, w int) string {
	if x < 0 {
		x = 0
	}
	if x > w {
		x = w
	}
	// Sicherstellen, dass base Breite w hat
	rb := []rune(base)
	if len(rb) < w {
		base = base + strings.Repeat(" ", w-len(rb))
		rb = []rune(base)
	} else if len(rb) > w {
		base = string(rb[:w])
		rb = []rune(base)
	}
	segRunes := []rune(segment)
	// Kopieren
	for i := 0; i < len(segRunes) && (x+i) < w; i++ {
		rb[x+i] = segRunes[i]
	}
	return string(rb)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// flattenOutput bricht Zeilen aus src auf gegebene Breite um und liefert die logischen Zeilen für das Rendering.
func flattenOutput(src []string, width int) []string {
	if width <= 0 || len(src) == 0 {
		return nil
	}
	out := make([]string, 0, len(src))
	for _, line := range src {
		r := []rune(line)
		for len(r) > width {
			out = append(out, string(r[:width]))
			r = r[width:]
		}
		out = append(out, string(r))
	}
	return out
}

// recalcViewport berechnet Viewport-Größe und abhängige Dimensionen neu.
func recalcViewport(m Model) Model {
	w := m.width
	h := m.height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}
	innerW := w - 2
	innerH := h - 2
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}
	header := RenderTitleBar(innerW, m.currentPage)
	headerHeight := 1
	if header != "" {
		headerHeight = len(strings.Split(header, "\n"))
	}
	// Footer-Höhe basierend auf aktuellem Hilfe-Status bestimmen
	m.help.Width = innerW
	footerRaw := m.help.View(m.keys)
	footerHeight := len(strings.Split(footerRaw, "\n"))

	// 3 Zeilen unten reservieren: 1 Trennlinie + 2 Fortschrittszeilen
	avail := innerH - headerHeight - footerHeight - 3
	if avail < 0 {
		avail = 0
	}
	outputRows := avail - 3 // Leerzeile nach Header + Trennzeile + Leerzeile vor Footer
	if outputRows < 0 {
		outputRows = 0
	}
	// Statusspalte (1/4 der Innenbreite) plus Trennzeichen (1 Zeichen) falls Breite reicht
	statusW := innerW / 4
	if statusW < 12 { // Mindestbreite erzwingen damit Hauptbereich nicht kollabiert
		statusW = 0
	}
	mainW := innerW - statusW
	if statusW > 0 {
		mainW -= 1 // Trennzeichen
	}
	if mainW < 0 {
		mainW = 0
	}
	m.vp.Width = mainW
	m.vp.Height = outputRows
	return m
}

// View rendert die komplette Oberfläche.
func (m Model) View() string {
	// Sinnvolle Mindestabmessungen früh sicherstellen
	w := m.width
	h := m.height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}

	// Innere Abmessungen (wir lassen einen leeren 1-Zeichen Rand)
	innerW := w - 2
	innerH := h - 2
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	// Header (mehrzeilig) über eigene Komponente bei innerer Breite
	header := RenderTitleBar(innerW, m.currentPage)
	headerHeight := 1
	if header != "" {
		headerHeight = len(strings.Split(header, "\n"))
	}

	// Footer über bubbles/help (variable Höhe möglich)
	m.help.Width = innerW
	footerRaw := m.help.View(m.keys)
	footer := footerLineStyle.Width(innerW).Render(footerRaw)
	footerHeight := len(strings.Split(footerRaw, "\n"))

	// Mittelteil via Dashboard/PullRequest Renderer zusammensetzen (benötigt Footerhöhe)
	// Wir rendern Inhalte auf einer effektiven Innenhöhe abzüglich der 2 Kachelzeilen
	effectiveInnerH := innerH - 3                      // Trennlinie + 2 Tile-Zeilen
	if effectiveInnerH < headerHeight+footerHeight+3 { // Mindeststruktur beibehalten
		effectiveInnerH = headerHeight + footerHeight + 3
	}
	// Haupt-Dashboard für reduzierte Breite (ohne Statusspalte) rendern und Zeilen anschließend zusammenführen
	statusW := innerW / 4
	if statusW < 12 {
		statusW = 0
	}
	mainW := innerW - statusW
	if statusW > 0 {
		mainW -= 1 // Trennzeichen
	}
	if mainW < 0 {
		mainW = 0
	}
	// Hilfe-Breite temporär für Berechnungen des Hauptbereichs anpassen
	origHelpWidth := m.help.Width
	m.help.Width = mainW
	var middleMain []string
	if m.currentPage == "Pull Requests" {
		middleMain = RenderPullRequests(m, mainW, effectiveInnerH, headerHeight, footerHeight)
	} else {
		middleMain = RenderDashboard(m, mainW, effectiveInnerH, headerHeight, footerHeight)
	}
	m.help.Width = origHelpWidth
	// Statusspalten-Zeilen erzeugen (gleiche Anzahl wie middleMain) falls aktiv
	var statusLines []string
	if statusW > 0 {
		statusLines = renderStatusColumn(m, statusW, len(middleMain))
	}
	// Spalten zeilenweise zusammenführen
	middle := make([]string, len(middleMain))
	for i, ln := range middleMain {
		if statusW > 0 {
			divider := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render("│")
			middle[i] = ln + divider + statusLines[i]
		} else {
			middle[i] = ln
		}
	}

	// Finale Zeilen mit genauer Höhe bauen
	lines := make([]string, 0, innerH)
	lines = append(lines, strings.Split(header, "\n")...)
	lines = append(lines, middle...)
	// Erst Footer (Hilfebar) anhängen
	lines = append(lines, footer)
	// Trennlinie über volle Breite
	divLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render(strings.Repeat("─", innerW))
	lines = append(lines, divLine)
	// Dann Fortschritts-Kacheln unter der Trennlinie über gesamte Breite (inkl. Statusbereich)
	progressLines := RenderProgressTiles(m, innerW)
	lines = append(lines, progressLines...)

	// Entfernte frühere zusätzliche Progress Tiles

	// Mit einem ein Zeichen breiten leeren Rand umschließen
	out := make([]string, 0, h)
	// Obere leere Randzeile
	out = append(out, strings.Repeat(" ", w))
	// Mittlere Zeilen mit einem Space links und rechts
	for _, ln := range lines {
		out = append(out, " "+ln+" ")
	}
	// Falls innerH == 0, Form beibehalten
	if innerH == 0 {
		out = append(out, strings.Repeat(" ", w))
	}
	// Untere leere Randzeile
	out = append(out, strings.Repeat(" ", w))
	return strings.Join(out, "\n")
}
