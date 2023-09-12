package task

import (
	"log"
	"testing"
)

func TestTaskAddTask(t *testing.T) {
	// queue a tasks that adds another task to current frame; should execute
	f := 1
	h := New(&f)

	count := 0
	h.Add(func() {
		log.Println("hello i'm at first task")
		count++

		h.Add(func() {
			count++
			log.Println("hello this is the second task")
		}, 0)

	}, 0)

	h.Run()

	if count != 2 {
		log.Printf("expecting count to be 2, got %v\n", count)
		t.FailNow()
	}

}
