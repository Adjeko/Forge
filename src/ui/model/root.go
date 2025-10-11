package model

import (
	"context"
	"fmt"
	"forge/src/exec"
	"forge/src/logging"
	"forge/src/output"
	"forge/src/ui/accessibility"
	"forge/src/ui/components"
	"forge/src/ui/help"
	"forge/src/ui/zones"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	cmdList         *components.CommandList
	viewport        *components.OutputViewport
	buffer          *output.OutputBuffer
	exitCode        int
	workflow        *exec.Workflow
	progressPercent int
	monitorsPanel   *components.MonitorsPanel
	helpOverlay     *help.Overlay
	focus           *accessibility.FocusManager
}

func NewRootModel() RootModel {
	buf := output.NewBuffer(100000)
	cmdList := components.NewCommandList(exec.Whitelist, buf)
	viewport := components.NewOutputViewport(buf, 50)
	// configurable workflow steps via env FORGE_WORKFLOW_STEPS (default 1 for perf tests)
	steps := 1
	if v := os.Getenv("FORGE_WORKFLOW_STEPS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			steps = n
		}
	}
	wfSteps := make([]exec.ExecutionStep, steps)
	for i := 0; i < steps; i++ {
		wfSteps[i] = exec.ExecutionStep{Cmd: exec.Whitelist[0]}
	}
	wf := &exec.Workflow{Steps: wfSteps}
	// initialize a few fake monitors (US3 placeholder)
	monitors := []exec.Monitor{
		exec.NewFakeMonitor("ping-local", 5*time.Second),
		exec.NewFakeMonitorWithType("script-daily", 5*time.Second, "script"),
	}
	panel := components.NewMonitorsPanel(monitors)
	// start scheduler (non-blocking) - could be moved to Init when more complex
	scheduler := exec.NewMonitorScheduler(monitors, 5*time.Second)
	go scheduler.Start()
	helpOverlay := help.NewOverlay()
	// Register help action
	accessibility.Register(accessibility.HotkeyAction{ID: "help", Description: "Toggle help overlay", Keys: []string{"?"}, ZoneID: "help:overlay"})
	zones.RegisterZone("help:overlay")
	// Register run workflow action parity
	accessibility.Register(accessibility.HotkeyAction{ID: "run-workflow", Description: "Run workflow", Keys: []string{"enter"}, ZoneID: "action:run-workflow"})
	zones.RegisterZone("action:run-workflow")
	zones.RegisterZone("action:scroll-up")
	zones.RegisterZone("action:scroll-down")
	focus := accessibility.NewFocusManager([]string{"CommandList", "Output", "Monitors", "Help"})
	return RootModel{cmdList: cmdList, viewport: viewport, buffer: buf, workflow: wf, monitorsPanel: panel, helpOverlay: helpOverlay, focus: focus}
}

func (m RootModel) Init() tea.Cmd {
	// Parity verification routine (T052): ensure each registered action has a zone.
	return func() tea.Msg {
		actions := accessibility.List()
		zoneSet := map[string]bool{}
		for _, z := range zones.Zones() {
			zoneSet[z] = true
		}
		for _, a := range actions {
			if !zoneSet[a.ZoneID] {
				m.buffer.Append("PARITY MISSING ZONE for action: " + a.ID)
			}
		}
		return nil
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "enter":
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			failedAt, err := exec.RunWorkflow(ctx, m.workflow, m.buffer)
			cancel()
			if err != nil {
				m.exitCode = -1
				if failedAt >= 0 {
					m.progressPercent = (failedAt * 100) / len(m.workflow.Steps)
				}
			} else {
				m.exitCode = 0
				m.progressPercent = 100
			}
		case "?":
			if m.helpOverlay != nil {
				m.helpOverlay.Toggle()
				logging.LogEvent(logging.EventUIHelpToggle, "visible", m.helpOverlay.Visible)
			}
		case "up":
			if m.cmdList.Selected > 0 {
				m.cmdList.Selected--
			}
		case "down":
			if m.cmdList.Selected < len(m.cmdList.Commands)-1 {
				m.cmdList.Selected++
			}
		case "tab":
			if m.focus != nil {
				m.focus.Next()
				logging.LogEvent(logging.EventUIFocusChange, "current", m.focus.Current())
			}
		case "pgup":
			if m.viewport != nil && m.focus != nil && m.focus.Current() == "Output" {
				m.viewport.ScrollUp()
				logging.LogEvent("ui.scroll", "dir", "up")
			}
		case "pgdown":
			if m.viewport != nil && m.focus != nil && m.focus.Current() == "Output" {
				m.viewport.ScrollDown()
				logging.LogEvent("ui.scroll", "dir", "down")
			}
		}
	}
	return m, nil
}

func (m RootModel) View() string {
	border := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	leftContent := "<commands unavailable>"
	if m.cmdList != nil {
		leftContent = m.cmdList.View()
	}
	rightContent := "<output unavailable>"
	if m.viewport != nil {
		rightContent = m.viewport.View()
		zones.RegisterZone("output:viewport")
	}
	focusStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	left := border.Render(leftContent)
	right := border.Render(rightContent)
	if m.focus != nil {
		current := m.focus.Current()
		if current == "CommandList" {
			left = focusStyle.Render(left)
		} else if current == "Output" {
			right = focusStyle.Render(right)
		}
	}
	main := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	bar := components.ProgressBar(m.progressPercent)
	monitorsView := ""
	if m.monitorsPanel != nil {
		mv := m.monitorsPanel.View()
		if m.focus != nil && m.focus.Current() == "Monitors" {
			mv = focusStyle.Render(mv)
		}
		monitorsView = " [monitors " + mv + "]"
	}
	footer := components.Footer() + " " + bar + monitorsView + lipgloss.NewStyle().Faint(true).Render(" exit="+fmtStatus(m.exitCode))
	helpView := ""
	if m.helpOverlay != nil {
		helpView = m.helpOverlay.View()
		if helpView != "" && m.focus != nil && m.focus.Current() == "Help" {
			helpView = focusStyle.Render(helpView)
		}
	}
	if helpView != "" {
		return components.Header() + "\n" + main + "\n" + helpView + "\n" + footer
	}
	return components.Header() + "\n" + main + "\n" + footer
}

func fmtStatus(code int) string {
	return fmt.Sprintf("%d", code)
}

// FocusCurrent exposes current focus region (test helper)
func (m RootModel) FocusCurrent() string {
	if m.focus == nil {
		return ""
	}
	return m.focus.Current()
}

// Buffer accessor for tests
func (m RootModel) Buffer() *output.OutputBuffer { return m.buffer }

// WorkflowLen returns number of steps (test helper for env flag)
func (m RootModel) WorkflowLen() int {
	if m.workflow == nil {
		return 0
	}
	return len(m.workflow.Steps)
}
