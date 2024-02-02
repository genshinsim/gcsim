// package queue provide a universal way of handling queuing and executing tasks
package queue

// TODO: This entire Task queue should be replaced and core's task queue should be used instead
// In order to replace, the core task queue must support the ability to update the position of a
// task in the queue.
// Also will need to consider order. Currently everything that is queued via QueueCharTask will
// always happen before all entries in the core task queue. If any implementations depend on this order,
// this will cause additional problems.
type Task struct {
	F     func()
	Delay int
}

func Add(slice *[]Task, f func(), delay int) {
	if delay == 0 {
		f()
		return
	}
	*slice = append(*slice, Task{
		F:     f,
		Delay: delay,
	})
}

func Run(slice *[]Task, currentTime int) {
	n := 0
	for i := 0; i < len(*slice); i++ {
		if (*slice)[i].Delay <= currentTime {
			(*slice)[i].F()
		} else {
			(*slice)[n] = (*slice)[i]
			n++
		}
	}
	*slice = (*slice)[:n]
}
