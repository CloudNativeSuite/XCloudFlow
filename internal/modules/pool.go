package modules

import "sync"

// Pool limits the number of concurrent goroutines.
type Pool struct {
	limit chan struct{}
	wg    sync.WaitGroup
}

// NewPool creates a pool that allows up to n concurrent tasks.
func NewPool(n int) *Pool {
	if n <= 0 {
		n = 1
	}
	return &Pool{limit: make(chan struct{}, n)}
}

// Submit runs fn in a new goroutine with pool concurrency control.
func (p *Pool) Submit(fn func()) {
	p.limit <- struct{}{}
	p.wg.Add(1)
	go func() {
		defer func() {
			<-p.limit
			p.wg.Done()
		}()
		fn()
	}()
}

// Wait blocks until all submitted tasks are done.
func (p *Pool) Wait() {
	p.wg.Wait()
}
