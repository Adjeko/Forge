package main

import (
	"bufio"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Streaming-Nachrichten für zeilenweise Befehlsausgabe.
type streamLineMsg struct {
	line  string
	isErr bool
}
type streamDoneMsg struct {
	err   string
	isErr bool
}

// readNextLine blockiert bis zur nächsten Zeile vom ausgewählten Scanner und liefert sie als Message zurück.
func (m Model) readNextLine(isErr bool) tea.Cmd {
	var scan *bufio.Scanner
	if isErr {
		scan = m.scannerErr
	} else {
		scan = m.scannerOut
	}
	return func() tea.Msg {
		if scan == nil {
			return streamDoneMsg{err: "scanner is nil", isErr: isErr}
		}
		if scan.Scan() {
			// Eventuelles CR unter Windows normalisieren
			line := scan.Text()
			line = strings.TrimSuffix(line, "\r")
			return streamLineMsg{line: line, isErr: isErr}
		}
		if err := scan.Err(); err != nil {
			return streamDoneMsg{err: err.Error(), isErr: isErr}
		}
		return streamDoneMsg{isErr: isErr}
	}
}

// StartStreamingCommand startet den angegebenen Befehl in Verzeichnis dir und initialisiert Streaming-Zustand am Model.
// Gibt ein Cmd zurück, das stdout und stderr gleichzeitig zu lesen beginnt.
func StartStreamingCommand(m *Model, name string, args []string, dir string) (tea.Cmd, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	m.streaming = true
	m.doneOut = false
	m.doneErr = false
	m.cmd = cmd
	scOut := bufio.NewScanner(stdout)
	scErr := bufio.NewScanner(stderr)
	// Lange Zeilen bis 1MB erlauben
	bufOut := make([]byte, 0, 64*1024)
	bufErr := make([]byte, 0, 64*1024)
	scOut.Buffer(bufOut, 1024*1024)
	scErr.Buffer(bufErr, 1024*1024)
	m.scannerOut = scOut
	m.scannerErr = scErr
	return tea.Batch(m.readNextLine(false), m.readNextLine(true)), nil
}
