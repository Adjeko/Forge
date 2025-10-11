package output

// OutputBuffer retains lines and issues advisory warning over a threshold.
type OutputBuffer struct {
	lines         []string // placeholder; future: []StyledLine
	warningIssued bool
	threshold     int
	advisoryLine  string // cached advisory message
}

// StyledLine placeholder struct for future styling tokens.
type StyledLine struct {
	Text string
}

func NewBuffer(threshold int) *OutputBuffer { return &OutputBuffer{threshold: threshold} }

// Append adds a line to the buffer retaining all history. When the configured threshold
// is exceeded for the first time an advisory line is appended exactly once to inform
// users of potential performance degradation. A tiny custom itoa avoids fmt allocations
// in this hot path under high-volume streaming tests (T089/T079).
func (b *OutputBuffer) Append(line string) {
	b.lines = append(b.lines, line)
	if !b.warningIssued && b.threshold > 0 && len(b.lines) > b.threshold {
		b.warningIssued = true
		// append advisory line once
		b.advisoryLine = "ADVISORY: output line threshold exceeded (" + itoa(len(b.lines)) + ")"
		b.lines = append(b.lines, b.advisoryLine)
	}
}

// itoa minimal local integer to string conversion (avoid fmt for micro-hot path)
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	b := [20]byte{}
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}

func (b *OutputBuffer) Len() int            { return len(b.lines) }
func (b *OutputBuffer) WarningIssued() bool { return b.warningIssued }

// ViewportSlice returns slice representing last n lines (retention behavior).
func (b *OutputBuffer) ViewportSlice(n int) []string {
	if n <= 0 || len(b.lines) <= n {
		return b.lines
	}
	return b.lines[len(b.lines)-n:]
}
