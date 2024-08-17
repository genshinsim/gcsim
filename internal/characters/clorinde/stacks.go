package clorinde

type queuer func(cb func(), delay int)

type stackTracker struct {
	queuer
	stacks []*int
	max    int
	frame  *int
}

func newStackTracker(maxstacks int, queue queuer, frame *int) *stackTracker {
	s := &stackTracker{
		queuer: queue,
		stacks: make([]*int, maxstacks),
		max:    maxstacks,
		frame:  frame,
	}
	return s
}

func (s *stackTracker) Count() int {
	count := 0
	for _, v := range s.stacks {
		if v != nil {
			count++
		}
	}
	return count
}

func (s *stackTracker) Add(duration int) {
	idx := 0
	for i := 0; i < s.max; i++ {
		if s.stacks[i] == nil {
			idx = i
			break
		}
		if *s.stacks[i] < *s.stacks[idx] {
			idx = i
		}
	}

	src := *s.frame
	s.stacks[idx] = &src

	s.queuer(func() {
		if s.stacks[idx] != &src {
			return
		}
		s.stacks[idx] = nil
	}, duration)
}
