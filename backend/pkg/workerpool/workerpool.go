package workerpool

import (
	"context"
	"errors"
	"sync"
)

type Task func(ctx context.Context)

type Pool struct {
	tasks  chan Task
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(size int) *Pool {
	if size <= 0 {
		size = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		tasks:  make(chan Task, size*2),
		ctx:    ctx,
		cancel: cancel,
	}

	p.wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-p.ctx.Done():
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					if task != nil {
						task(p.ctx)
					}
				}
			}
		}()
	}

	return p
}

func (p *Pool) Submit(task Task) error {
	if task == nil {
		return nil
	}

	select {
	case <-p.ctx.Done():
		return errors.New("workerpool: pool is stopped")
	default:
	}

	select {
	case p.tasks <- task:
		return nil
	case <-p.ctx.Done():
		return errors.New("workerpool: pool is stopped")
	}
}

func (p *Pool) Stop() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
}
