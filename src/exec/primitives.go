package exec

// PrimitiveCommand represents a whitelisted command.
type PrimitiveCommand struct {
	ID    string
	Label string
	Cmd   string
}

// Whitelist contains allowed primitive commands (stub placeholder).
var Whitelist = []PrimitiveCommand{
	{ID: "echo-ok", Label: "echo ok", Cmd: "cmd /C echo ok"},
}

// IsWhitelisted returns true if command string exactly matches a whitelisted command.
func IsWhitelisted(cmd string) bool {
	for _, p := range Whitelist {
		if p.Cmd == cmd {
			return true
		}
	}
	return false
}
