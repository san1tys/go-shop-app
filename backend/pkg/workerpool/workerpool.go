package workerpool

import (
	"context"
	"errors"
	"sync"
)

// Task описывает задачу, выполняемую воркером.
type Task func(ctx context.Context)

// Pool — простой worker pool поверх горутин и канала задач.
type Pool struct {
	tasks  chan Task
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New создаёт и запускает пул из size воркеров.
func New(size int) *Pool {
	if size <= 0 {
		size = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		// небольшой буфер, чтобы не блокировать короткие всплески задач
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

// Submit отправляет задачу в пул.
// Возвращает ошибку, если пул уже остановлен.
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

// Stop останавливает пул: перестаёт принимать задачи и ждёт завершения воркеров.
func (p *Pool) Stop() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
}
