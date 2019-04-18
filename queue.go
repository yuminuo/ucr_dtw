package UcrDtw

type Queue struct {
	items []float64
}

func (s *Queue) Enqueue(t float64) {
	s.items = append(s.items, t)
}

func (s *Queue) Dequeue() float64 {
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	return item
}

func (s *Queue) Count() int {
	return len(s.items)
}
