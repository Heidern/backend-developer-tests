// Package concurrency implements worker pool interfaces, one simple and one a
// bit more complex.
package concurrency

// SimplePool is a simple worker pool that does not support cancellation or
// closing. All functions are safe to call from multiple goroutines.
type SimplePool interface {
	// Submit a task to be executed asynchronously. This function will return as
	// soon as the task is submitted. If the pool does not have an available slot
	// for the task, this blocks until it can submit.
	Submit(func())
}

type DefaultSimplePool struct {
	maxConcurrent     int
	countingSemaphore chan struct{}
}

// NewSimplePool creates a new SimplePool that only allows the given maximum
// concurrent tasks to run at any one time. maxConcurrent must be greater than
// zero.
func NewSimplePool(maxConcurrent int) SimplePool {
	return &DefaultSimplePool{
		maxConcurrent:     maxConcurrent,
		countingSemaphore: make(chan struct{}, maxConcurrent),
	}
}

func (pool *DefaultSimplePool) Submit(task func()) {
	pool.countingSemaphore <- struct{}{}
	go func() {
		task()
		<-pool.countingSemaphore
	}()
}
