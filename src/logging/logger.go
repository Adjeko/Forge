package logging

import (
	"log/slog"
	"time"
)

var Logger = slog.Default()

const (
	EventCommandStart      = "command.start"
	EventCommandEnd        = "command.end"
	EventWorkflowStart     = "workflow.start"
	EventWorkflowEnd       = "workflow.end"
	EventMonitorPollStart  = "monitor.poll.start"
	EventMonitorPollResult = "monitor.poll.result"
	EventUIFocusChange     = "ui.focus.change"
	EventUIHelpToggle      = "ui.help.toggle"
)

type Event struct {
	Name      string
	Timestamp time.Time
	Attrs     []any
}

var captured []Event

func LogEvent(event string, attrs ...any) {
	Logger.Info(event, attrs...)
	captured = append(captured, Event{Name: event, Timestamp: time.Now(), Attrs: attrs})
}

func CapturedEvents() []Event { return append([]Event{}, captured...) }
func ResetEvents()            { captured = nil }
