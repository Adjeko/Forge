package main

import (
	"fmt"
	"forge/src/ui/model"

	tea "github.com/charmbracelet/bubbletea"
)

type quitWrapper struct{ inner model.RootModel }

func (m quitWrapper) Init() tea.Cmd { return m.inner.Init() }
func (m quitWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok && k.String() == "q" {
		return m, tea.Quit
	}
	_, cmd := m.inner.Update(msg)
	return m, cmd
}
func (m quitWrapper) View() string { return m.inner.View() }

func main() {
	p := tea.NewProgram(quitWrapper{inner: model.NewRootModel()})
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
	}
}
