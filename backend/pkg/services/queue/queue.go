package queue

import (
	"errors"
	"sync"
	"time"
)

type Queue struct {
	Timeout int
	work    []*Work
	wip     map[string]*Work
	mu      sync.Mutex
}

type Work struct {
	Key  string
	Work any
}

func NewQueue(timeout int) *Queue {
	q := &Queue{
		Timeout: timeout,
		work:    make([]*Work, 0, 10),
	}
	if q.Timeout <= 0 {
		q.Timeout = 5 * 60 //default 5 min
	}
	return q
}

func (q *Queue) Complete(key string) error {
	//removes the work from wip if exists
	//otherwise errors
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.wip[key]; !ok {
		return errors.New("key is not in wip")
	}
	delete(q.wip, key)
	return nil
}

func (q *Queue) Pop() *Work {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.work) == 0 {
		return nil
	}
	x := q.work[0]
	q.work = q.work[1:]
	//add this to wip and implement a 5min timeout
	go q.deleteAfterTimeout(x.Key)
	return x
}

func (q *Queue) deleteAfterTimeout(key string) {
	time.Sleep(time.Duration(q.Timeout) * time.Second)
	q.mu.Lock()
	defer q.mu.Unlock()
	w, ok := q.wip[key]
	if !ok {
		return
	}
	q.Insert(w)
	delete(q.wip, key)
}

func (q *Queue) Append(w *Work) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.work = append(q.work, w)
}

func (q *Queue) Insert(w *Work) {
	q.mu.Lock()
	defer q.mu.Unlock()
	//we assume key is unique... if not then things will probably break
	q.work = append([]*Work{w}, q.work...)
}
