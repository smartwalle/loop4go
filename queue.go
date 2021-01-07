package loop4go

import "sync"

// Queue https://github.com/davyxu/cellnet/blob/master/pipe.go
type Queue struct {
	items []interface{}
	mu    sync.Mutex
	cond  *sync.Cond
}

func NewQueue() *Queue {
	var q = &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (this *Queue) Enqueue(item interface{}) {
	this.mu.Lock()
	this.items = append(this.items, item)
	this.mu.Unlock()
	this.cond.Signal()
}

func (this *Queue) Reset() {
	this.items = this.items[0:0]
}

func (this *Queue) Dequeue(items *[]interface{}) (exit bool) {
	this.mu.Lock()
	for len(this.items) == 0 {
		this.cond.Wait()
	}
	this.mu.Unlock()

	this.mu.Lock()
	for _, item := range this.items {
		if item == nil {
			exit = true
			break
		} else {
			*items = append(*items, item)
		}
	}

	this.Reset()
	this.mu.Unlock()
	return exit
}
