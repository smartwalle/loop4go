package loop4go

import (
	"sync/atomic"
	"time"
)

type Loop struct {
	duration time.Duration
	running  int32
	queue    EventQueue
	callback func()
}

func NewLoop(d time.Duration, queue EventQueue, callback func()) *Loop {
	var t = &Loop{}
	t.duration = d
	t.running = 0
	t.queue = queue
	t.callback = callback
	return t
}

func (this *Loop) Running() bool {
	return atomic.LoadInt32(&this.running) == 1
}

func (this *Loop) Start() bool {
	if this.duration <= 0 {
		return false
	}

	if old := atomic.SwapInt32(&this.running, 1); old != 0 {
		return false
	}

	this.enqueue()

	return true
}
func (this *Loop) Stop() {
	if old := atomic.SwapInt32(&this.running, 0); old != 1 {
		return
	}
}

func (this *Loop) enqueue() {
	if this.Running() {
		after(this.duration, this.queue, func() {
			this.exec()
		})
	}
}

func (this *Loop) exec() {
	if this.Running() {
		defer this.enqueue()
	}
	this.callback()
}

func after(d time.Duration, q EventQueue, callback func()) {
	time.AfterFunc(d, func() {
		if q != nil {
			q.Enqueue(callback)
		} else {
			callback()
		}
	})
}
