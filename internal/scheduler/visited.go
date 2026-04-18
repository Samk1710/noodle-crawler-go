package scheduler

import "sync"

type Visited struct {
	mu   sync.Mutex
	data map[string]bool
}

func NewVisited() *Visited {
	return &Visited{
		data: make(map[string]bool),
	}
}

func (v *Visited) Add(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.data[url] {
		return false
	}

	v.data[url] = true
	return true
}