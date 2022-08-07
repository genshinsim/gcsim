//package queue provide a universal way of handling queuing and executing tasks
package queue

type Task struct {
	F     func()
	Delay float64
}

func Add(slice *[]Task, f func(), delay float64) {
	if delay == 0 {
		f()
		return
	}
	*slice = append(*slice, Task{
		F:     f,
		Delay: delay,
	})
}

func Run(slice *[]Task, currentTime float64) {
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
