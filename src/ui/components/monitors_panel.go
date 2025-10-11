package components

import (
	"forge/src/exec"
	"forge/src/ui/zones"
	"strings"
)

// MonitorsPanel renders LED indicators for current monitors.
type MonitorsPanel struct {
	Monitors []exec.Monitor
}

func NewMonitorsPanel(monitors []exec.Monitor) *MonitorsPanel {
	return &MonitorsPanel{Monitors: monitors}
}

// View returns a simple horizontal list of LEDs: [ID:color]
func (mp *MonitorsPanel) View() string {
	if mp == nil || len(mp.Monitors) == 0 {
		return "(no monitors)"
	}
	parts := make([]string, 0, len(mp.Monitors))
	for _, m := range mp.Monitors {
		color := exec.MonitorLEDColor(m.State())
		// Register zone per monitor LED (placeholder; future click semantics)
		zones.RegisterZone("monitor:" + m.ID())
		parts = append(parts, m.ID()+":"+color)
	}
	return strings.Join(parts, " ")
}
