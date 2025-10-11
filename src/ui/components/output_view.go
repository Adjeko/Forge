package components

import (
	"forge/src/output"
	"strings"
)

type OutputViewport struct {
	Buf         *output.OutputBuffer
	Window      int
	Offset      int    // starting line index from tail (0 = most recent)
	cacheJoined string // cached joined representation for current tail slice when Offset==0
	cacheLen    int    // number of lines included in cache
}

func NewOutputViewport(buf *output.OutputBuffer, window int) *OutputViewport {
	return &OutputViewport{Buf: buf, Window: window}
}

func (v *OutputViewport) View() string {
	lines := v.Buf.ViewportSlice(v.Window)
	if v.Offset == 0 {
		// Use cache if length matches
		if v.cacheLen == len(lines) && v.cacheJoined != "" {
			return v.cacheJoined
		}
		joined := strings.Join(lines, "\n")
		v.cacheJoined = joined
		v.cacheLen = len(lines)
		return joined
		// Scrolling invalidates the tail cache and recomputes a trimmed slice.
	}
	// simple offset from end: if Offset > 0, display an earlier slice placeholder
	if v.Offset > 0 && v.Offset < len(lines) {
		lines = lines[:len(lines)-v.Offset]
	}
	return strings.Join(lines, "\n")
}

func (v *OutputViewport) ScrollUp() {
	if v.Offset+1 < v.Window { // limit simplistic for prototype
		v.Offset++
	}
}

func (v *OutputViewport) ScrollDown() {
	if v.Offset > 0 {
		v.Offset--
	}
	if v.Offset == 0 { // invalidate cache so tail refresh reflects changes
		v.cacheJoined = ""
		v.cacheLen = 0
	}
}

// SetWindow adjusts the viewport window size and invalidates cache.
func (v *OutputViewport) SetWindow(n int) {
	if n <= 0 {
		return
	}
	v.Window = n
	v.cacheJoined = ""
	v.cacheLen = 0
	if v.Offset >= v.Window { // clamp offset
		v.Offset = v.Window - 1
		if v.Offset < 0 {
			v.Offset = 0
		}
	}
}
