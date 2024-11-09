package stacks

type queuer func(cb func(), delay int)

type MultipleRefreshNoRemove struct {
	queuer
	stacks []*int
	frame  *int
}

func NewMultipleRefreshNoRemove(maxstacks int, queue queuer, frame *int) *MultipleRefreshNoRemove {
	s := &MultipleRefreshNoRemove{
		queuer: queue,
		stacks: make([]*int, maxstacks),
		frame:  frame,
	}
	return s
}

func (s *MultipleRefreshNoRemove) Count() int {
	count := 0
	for _, v := range s.stacks {
		if v != nil {
			count++
		}
	}
	return count
}

func (s *MultipleRefreshNoRemove) Add(duration int) {
	idx := 0
	for i := 0; i < len(s.stacks); i++ {
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
