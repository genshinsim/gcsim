package combat

func (s *Sim) AddStatus(key string, dur int) {
	s.status[key] = s.f + dur
}

func (s *Sim) DeleteStatus(key string) {
	delete(s.status, key)
}

func (s *Sim) Status(key string) int {
	f, ok := s.status[key]
	if !ok {
		return 0
	}
	if f > s.f {
		return f - s.f
	}
	return 0
}
