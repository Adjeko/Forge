package accessibility

// FocusManager manages cyclic traversal of focusable UI regions.
type FocusManager struct {
	regions []string
	index   int
}

func NewFocusManager(regions []string) *FocusManager { return &FocusManager{regions: regions} }

func (f *FocusManager) Current() string {
	if len(f.regions) == 0 {
		return ""
	}
	return f.regions[f.index]
}

func (f *FocusManager) Next() {
	if len(f.regions) > 0 {
		f.index = (f.index + 1) % len(f.regions)
	}
}

func (f *FocusManager) Set(name string) {
	for i, r := range f.regions {
		if r == name {
			f.index = i
			return
		}
	}
}

func (f *FocusManager) Regions() []string { return append([]string{}, f.regions...) }
