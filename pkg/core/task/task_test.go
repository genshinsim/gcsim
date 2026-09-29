package task

import (
	"container/heap"
	"log"
	"math/rand/v2"
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

// refHeap is the container/heap implementation that minHeap replaced.
type refHeap []task

func (h refHeap) Len() int { return len(h) }
func (h refHeap) Less(i, j int) bool {
	return h[i].executeBy < h[j].executeBy || (h[i].executeBy == h[j].executeBy && h[i].id < h[j].id)
}
func (h refHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *refHeap) Push(x any)   { *h = append(*h, x.(task)) }
func (h *refHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TestHeapMatchesContainerHeap(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	var got minHeap
	var want refHeap
	id := 0
	for step := range 50000 {
		if len(want) > 0 && rng.IntN(100) < 48 {
			g, w := got.pop(), heap.Pop(&want).(task)
			if g.executeBy != w.executeBy || g.id != w.id {
				t.Fatalf("step %v: pop = (%v, %v), want (%v, %v)", step, g.executeBy, g.id, w.executeBy, w.id)
			}
		} else {
			// few distinct frames so that ties on executeBy are common
			x := task{executeBy: rng.IntN(40), id: id}
			id++
			got.push(x)
			heap.Push(&want, x)
		}
		if len(got) != len(want) {
			t.Fatalf("step %v: len = %v, want %v", step, len(got), len(want))
		}
		for i := range got {
			if got[i].executeBy != want[i].executeBy || got[i].id != want[i].id {
				t.Fatalf("step %v: layout differs at %v", step, i)
			}
		}
	}
}
