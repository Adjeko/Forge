package accessibility

// HotkeyAction describes a dual-modality action.
type HotkeyAction struct {
	ID          string
	Description string
	Keys        []string
	ZoneID      string
}

var registry []HotkeyAction

func Register(a HotkeyAction) { registry = append(registry, a) }
func List() []HotkeyAction    { return registry }

// Reset clears the registry (test helper to avoid state leakage across tests).
func Reset() { registry = nil }

// ParityOK checks hotkey vs zone parity (zoneID present when keys present).
func ParityOK() bool {
	for _, a := range registry {
		if len(a.Keys) == 0 || a.ZoneID == "" {
			return false
		}
	}
	return true
}
