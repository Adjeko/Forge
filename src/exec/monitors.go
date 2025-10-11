package exec

import (
	"forge/src/logging"
	"math/rand"
	"sync"
	"time"
)

// MonitorState represents lifecycle/result states for a monitor.
type MonitorState int

const (
	StateChecking MonitorState = iota
	StateOK
	StateFailed
	StateHanging
	StateDisabled
)

// Monitor interface for polling background resources.
type Monitor interface {
	ID() string
	Interval() time.Duration
	Poll() // performs a check, updates internal state
	State() MonitorState
}

// FakeMonitor used in tests; simulates polling activity.
type FakeMonitor struct {
	id          string
	interval    time.Duration
	mu          sync.Mutex
	state       MonitorState
	pollCount   int
	lastPoll    time.Time
	monitorType string
	scriptPath  string
}

func NewFakeMonitor(id string, interval time.Duration) *FakeMonitor {
	return &FakeMonitor{id: id, interval: interval, state: StateChecking}
}

func NewFakeMonitorWithType(id string, interval time.Duration, typ string) *FakeMonitor {
	return &FakeMonitor{id: id, interval: interval, state: StateChecking, monitorType: typ}
}

func (f *FakeMonitor) ID() string              { return f.id }
func (f *FakeMonitor) Interval() time.Duration { return f.interval }
func (f *FakeMonitor) State() MonitorState {
	f.mu.Lock()
	defer f.mu.Unlock()
	// Hanging heuristic: no poll for >2*interval while in checking
	if f.state == StateChecking && !f.lastPoll.IsZero() && time.Since(f.lastPoll) > 2*f.interval {
		return StateHanging
	}
	return f.state
}
func (f *FakeMonitor) Poll() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pollCount++
	f.lastPoll = time.Now()
	// simple alternating success/failure pattern for simulation.
	if f.pollCount%2 == 0 {
		f.state = StateOK
	} else {
		f.state = StateFailed // purposely alternate for variance
	}
}

// Test helper methods
func (f *FakeMonitor) PollCount() int          { f.mu.Lock(); defer f.mu.Unlock(); return f.pollCount }
func (f *FakeMonitor) LastPollTime() time.Time { f.mu.Lock(); defer f.mu.Unlock(); return f.lastPoll }
func (f *FakeMonitor) SetResult(ok bool, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if ok {
		f.state = StateOK
	} else {
		f.state = StateFailed
	}
}
func (f *FakeMonitor) SetChecking()            { f.mu.Lock(); f.state = StateChecking; f.mu.Unlock() }
func (f *FakeMonitor) SetLastPoll(t time.Time) { f.mu.Lock(); f.lastPoll = t; f.mu.Unlock() }
func (f *FakeMonitor) SetScriptPath(p string)  { f.mu.Lock(); f.scriptPath = p; f.mu.Unlock() }
func (f *FakeMonitor) ScriptPath() string      { f.mu.Lock(); defer f.mu.Unlock(); return f.scriptPath }

// MonitorLEDColor maps a state to a LED color name (simplified string for tests).
func MonitorLEDColor(s MonitorState) string {
	switch s {
	case StateOK:
		return "green"
	case StateFailed:
		return "red"
	case StateHanging:
		return "yellow"
	case StateDisabled:
		return "gray"
	default:
		return "blue" // checking
	}
}

// MonitorScheduler coordinates polling across monitors with jitter.
type MonitorScheduler struct {
	monitors     []Monitor
	baseInterval time.Duration
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

func NewMonitorScheduler(monitors []Monitor, interval time.Duration) *MonitorScheduler {
	return &MonitorScheduler{monitors: monitors, baseInterval: interval, stopCh: make(chan struct{})}
}

func (s *MonitorScheduler) Start() {
	// naive implementation: launch goroutine per monitor.
	for _, m := range s.monitors {
		s.wg.Add(1)
		go func(m Monitor) {
			defer s.wg.Done()
			// initial jitter up to 50% interval
			jitter := time.Duration(rand.Int63n(int64(m.Interval()/2 + 1)))
			timer := time.NewTimer(jitter)
			for {
				select {
				case <-s.stopCh:
					timer.Stop()
					return
				case <-timer.C:
					logging.LogEvent(logging.EventMonitorPollStart, "id", m.ID())
					m.Poll()
					logging.LogEvent(logging.EventMonitorPollResult, "id", m.ID(), "state", MonitorLEDColor(m.State()))
					// schedule next poll with jitter (+0..50%)
					base := m.Interval()
					spread := int64(base / 2)
					delta := rand.Int63n(spread)
					next := base + time.Duration(delta)
					timer.Reset(next)
				}
			}
		}(m)
	}
}

func (s *MonitorScheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}
