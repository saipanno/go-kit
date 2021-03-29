package utils 

import "sync"

// WaitGroupWrapper ...
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap ...
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
