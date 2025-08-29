package skirk

import "fmt"

/// RingQueue modified from https://github.com/logdyhq/logdy-core/blob/main/ring/ring.go

type RingQueue[T any] struct {
	data   []T  // container data of a generic type T
	isFull bool // disambiguate whether the queue is full or empty
	start  int  // start index (inclusive, i.e. first element)
	end    int  // end index (exclusive, i.e. next after last element)
}

func NewRingQueue[T any](capacity int64) RingQueue[T] {
	return RingQueue[T]{
		data:   make([]T, capacity),
		isFull: false,
		start:  0,
		end:    0,
	}
}

func (r *RingQueue[T]) String() string {
	return fmt.Sprintf(
		"[RingQ full:%v size:%d start:%d end:%d data:%v]",
		r.isFull,
		len(r.data),
		r.start,
		r.end,
		r.data)
}

func (r *RingQueue[T]) Push(elem T) error {
	if r.isFull {
		return fmt.Errorf("out of bounds push, container is full")
	}
	r.pushUnchecked(elem)

	return nil
}

func (r *RingQueue[T]) PushOverwrite(elem T) {
	if r.isFull {
		r.Pop()
		r.pushUnchecked(elem)
		return
	}
	r.pushUnchecked(elem)
}

func (r *RingQueue[T]) pushUnchecked(elem T) {
	r.data[r.end] = elem              // place the new element on the available space
	r.end = (r.end + 1) % len(r.data) // move the end forward by modulo of capacity
	r.isFull = r.end == r.start       // check if we're full now
}

func (r *RingQueue[T]) Pop() (T, error) {
	var res T // "zero" element (respective of the type)
	if r.IsEmpty() {
		return res, fmt.Errorf("empty queue")
	}

	res = r.data[r.start]                 // copy over the first element in the queue
	r.start = (r.start + 1) % len(r.data) // move the start of the queue
	r.isFull = false                      // since we're removing elements, we can never be full

	return res, nil
}

func (r *RingQueue[T]) IsEmpty() bool {
	if !r.isFull && r.start == r.end {
		return true
	}
	return false
}

func (r *RingQueue[T]) IsFull() bool {
	return r.isFull
}

func (r *RingQueue[T]) Len() int {
	if r.isFull {
		return len(r.data)
	}
	return (r.end - r.start + len(r.data)) % len(r.data)
}

func (r *RingQueue[T]) Index(ind int) (T, error) {
	var res T // "zero" element (respective of the type)
	if ind >= r.Len() {
		return res, fmt.Errorf("Index out of bound")
	}
	return r.data[(r.start+ind)%len(r.data)], nil
}

func (r *RingQueue[T]) Clear() {
	r.start = r.end
	r.isFull = false
}

func (r *RingQueue[T]) Count(filter func(x T) bool) int {
	count := 0
	for i := range r.Len() {
		val, _ := r.Index(i)
		if filter(val) {
			count++
		}
	}
	return count
}
