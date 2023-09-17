package event

import (
	"sync"
	"time"
)

// Merge merges many tasks with same key into one task in a tiny interval.
type Merge interface {
	// UpdateOnline provides change time.Duration on running time.
	UpdateOnline(duration time.Duration)
	// Allowed tells your goroutine to go on or return directly.
	Allowed(key Key) bool
	// Clear clears local storage of task keys.
	Clear()
	// Group provides tree-Merge.
	Group(path string, newInterval time.Duration) Merge
	Run() error
	Close()
}

type void struct{}

var noop void

type Key string

type merge struct {
	storage  map[Key]void
	ticker   *time.Ticker
	mut      sync.RWMutex
	children map[string]Merge
}

var _ Merge = &merge{}

func NewMerge(interval time.Duration) Merge {
	return &merge{
		storage:  make(map[Key]void),
		ticker:   time.NewTicker(interval),
		children: make(map[string]Merge),
	}
}

func (m *merge) Group(path string, newInterval time.Duration) Merge {
	if me, has := m.children[path]; has {
		return me
	}
	child := NewMerge(newInterval)
	m.children[path] = child
	return child
}

func (m *merge) UpdateOnline(duration time.Duration) {
	if duration < time.Millisecond {
		return
	}
	m.ticker.Reset(duration)
}

func (m *merge) Allowed(key Key) bool {
	m.mut.RLock()
	_, has := m.storage[key]
	m.mut.RUnlock()
	if has {
		return false
	}

	m.mut.Lock()
	m.storage[key] = noop
	m.mut.Unlock()

	return true
}

func (m *merge) Clear() {
	var newStorage = make(map[Key]void)
	m.mut.Lock()
	m.storage = newStorage
	m.mut.Unlock()
}

func (m *merge) Run() error {
	go func() {
		for range m.ticker.C {
			m.Clear()
		}
	}()
	return nil
}

func (m *merge) Close() {
	m.ticker.Stop()
}
