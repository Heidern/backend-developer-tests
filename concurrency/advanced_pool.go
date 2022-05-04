package concurrency

import (
	"context"
	"errors"
	"sync"
)

// ErrPoolClosed is returned from AdvancedPool.Submit when the pool is closed
// before submission can be sent.
var ErrPoolClosed = errors.New("pool closed")

// AdvancedPool is a more advanced worker pool that supports cancelling the
// submission and closing the pool. All functions are safe to call from multiple
// goroutines.
type AdvancedPool interface {
	// Submit submits the given task to the pool, blocking until a slot becomes
	// available or the context is closed. The given context and its lifetime only
	// affects this function and is not the context passed to the callback. If the
	// context is closed before a slot becomes available, the context error is
	// returned. If the pool is closed before a slot becomes available,
	// ErrPoolClosed is returned. Otherwise the task is submitted to the pool and
	// no error is returned. The context passed to the callback will be closed
	// when the pool is closed.
	Submit(context.Context, func(context.Context)) error

	// Close closes the pool and waits until all submitted tasks have completed
	// before returning. If the pool is already closed, ErrPoolClosed is returned.
	// If the given context is closed before all tasks have finished, the context
	// error is returned. Otherwise, no error is returned.
	Close(context.Context) error
}

type DefaultAdvancedPool struct {
	maxSlots      int
	maxConcurrent int
	pendingTasks  chan func(context.Context)
	ctx           context.Context
	ctxCancel     context.CancelFunc
	taskGroup     sync.WaitGroup
}

// NewAdvancedPool creates a new AdvancedPool. maxSlots is the maximum total
// submitted tasks, running or waiting, that can be submitted before Submit
// blocks waiting for more room. maxConcurrent is the maximum tasks that can be
// running at any one time. An error is returned if maxSlots is less than
// maxConcurrent or if either value is not greater than zero.

func NewAdvancedPool(maxSlots, maxConcurrent int) (AdvancedPool, error) {
	poolContext, poolContextCancelFunc := context.WithCancel(context.Background())

	pool := &DefaultAdvancedPool{
		maxSlots:      maxSlots,
		maxConcurrent: maxConcurrent,
		pendingTasks:  make(chan func(context.Context), maxSlots-maxConcurrent),
		ctx:           poolContext,
		ctxCancel:     poolContextCancelFunc,
	}

	pool.taskGroup.Add(pool.maxConcurrent)

	for i := 0; i < pool.maxConcurrent; i++ {
		go func() {
			for t := range pool.pendingTasks {
				t(pool.ctx)
			}
			pool.taskGroup.Done()
		}()
	}

	return pool, nil
}

func (pool *DefaultAdvancedPool) Submit(ctx context.Context, task func(context.Context)) error {
	select {
	case <-pool.ctx.Done():
		return ErrPoolClosed

	case <-ctx.Done():
		return ctx.Err()

	case pool.pendingTasks <- task:
	}

	return nil
}

// Close closes the pool and waits until all submitted tasks have completed
// before returning.
// - If the pool is already closed, ErrPoolClosed is returned.
// - If the given context is closed before all tasks have finished, the context
//   error is returned.
// - Otherwise, no error is returned.

func (pool *DefaultAdvancedPool) Close(ctx context.Context) error {
	select {
	case <-pool.ctx.Done():
		return ErrPoolClosed
	default:
		pool.ctxCancel()
	}

	tasksFinished := make(chan struct{})

	go func() {
		close(pool.pendingTasks)

		pool.taskGroup.Wait()

		tasksFinished <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-tasksFinished:
		return nil
	}
}
