package zones

// RegisterZone placeholder for BubbleZone usage.
var registered []string

func RegisterZone(id string) { registered = append(registered, id) }
func Zones() []string        { return registered }

// Reset clears zones registry (test helper).
func Reset() { registered = nil }
