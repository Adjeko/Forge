package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

// Model-Definition und zugehörige Methoden ausgelagert in model.go

func main() {
	// Kontext mit Signal-Abbruch (Ctrl+C) zusätzlich zur Bubble Tea eigenen Behandlung.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Alt-Screen aktivieren um gesamten Terminalbereich exklusiv zu nutzen.
	p := tea.NewProgram(initialModel(), tea.WithContext(ctx), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		// Fehler wird klar ausgegeben, kein Panic.
		fmt.Fprintf(os.Stderr, "Fehler beim Ausführen der TUI: %v\n", err)
		os.Exit(1)
	}
}
