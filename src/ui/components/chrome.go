package components

import (
	"sync"
)

// Cached gradient-like header/footer (placeholder styling) to avoid recomputation under heavy renders.
var (
	headerOnce   sync.Once
	footerOnce   sync.Once
	cachedHeader string
	cachedFooter string
)

func Header() string {
	headerOnce.Do(func() {
		// Simulate gradient construction (cheap placeholder)
		cachedHeader = "FORGE TASK RUNNER"
	})
	return cachedHeader
}

func Footer() string {
	footerOnce.Do(func() {
		cachedFooter = "READY"
	})
	return cachedFooter
}
