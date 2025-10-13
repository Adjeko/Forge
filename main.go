package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		panic(err)
	}
}
