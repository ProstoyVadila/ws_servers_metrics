package gopool

// Goroutine pool manager package

import (
	"fmt"
	"log"
	"time"
)

// ErrScheduleTimeout returned by Pool to indicate that there no free
// goroutines during some period of time.
var ErrScheduleTimeout = fmt.Errorf("schedule error: timed out")

type Pool struct {
	sem  chan struct{}
	work chan func()
}

func New(size, queue, spawn int) *Pool {
	if spawn <= 0 && queue > 0 {
		log.Fatal("dead queue configuration detected")
	}
	if spawn > size {
		log.Fatal("spawn > workers")
	}
	p := &Pool{
		sem:  make(chan struct{}, size),
		work: make(chan func(), queue),
	}
	for i := 0; i < spawn; i++ {
		p.sem <- struct{}{}
		go p.worker(func() {})
	}
	return p
}

func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()

	task()
	for task := range p.work {
		task()
	}
}

// Schedule schedules task to be executed over pool's workers.
func (p *Pool) Schedule(task func()) {
	p.schedule(task, nil)
}

// ScheduleTimeout schedules task to be executed over pool's workers.
// It returns ErrScheduleTimeout when no free workers met during given timeout.
func (p *Pool) ScheduleTimeout(timeout time.Duration, task func()) error {
	return p.schedule(task, time.After(timeout))
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.work <- task:
		return nil
	case p.sem <- struct{}{}:
		go p.worker(task)
		return nil
	}
}
