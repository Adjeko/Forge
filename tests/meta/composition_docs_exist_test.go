package meta_test

import (
	"os"
	"path/filepath"
	"testing"
)

// T086G verify composition attempt docs exist
func TestCompositionDocsExist(t *testing.T) {
	cwd, _ := os.Getwd()
	root := cwd
	// attempt to locate repo root by ascending until specs exists (simplistic)
	for i := 0; i < 4; i++ {
		if _, err := os.Stat(filepath.Join(root, "docs", "composition")); err == nil {
			break
		}
		root = filepath.Dir(root)
	}
	files := []string{"commandlist.md", "outputviewport.md", "progressbar.md", "monitorspanel.md", "helpoverlay.md"}
	for _, f := range files {
		path := filepath.Join(root, "docs", "composition", f)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing composition doc: %s", path)
		}
	}
}
