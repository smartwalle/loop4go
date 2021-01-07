package loop4go

import (
	"sync"
)

// EventQueue https://github.com/davyxu/cellnet/blob/master/queue.go
type EventQueue interface {
	Start()

	Stop()

	Wait()

	Enqueue(func())
}

type eventQueue struct {
	q *Queue
	w *sync.WaitGroup
}

func NewEventQueue() EventQueue {
	var eq = &eventQueue{}
	eq.q = NewQueue()
	eq.w = &sync.WaitGroup{}
	return eq
}

func (this *eventQueue) Start() {
	this.w.Add(1)
	go func() {
		var itemList []interface{}
		var exit bool
		for {
			itemList = itemList[0:0]
			exit = this.q.Dequeue(&itemList)

			for _, item := range itemList {
				switch m := item.(type) {
				case func():
					m()
				case nil:
					break
				}
			}

			if exit {
				break
			}
		}
		this.w.Done()
	}()
}

func (this *eventQueue) Stop() {
	this.q.Enqueue(nil)
}

func (this *eventQueue) Wait() {
	this.w.Wait()
}

func (this *eventQueue) Enqueue(f func()) {
	if f == nil {
		return
	}
	this.q.Enqueue(f)
}
