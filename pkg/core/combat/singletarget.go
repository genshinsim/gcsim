package combat

import "fmt"

type SingleTarget struct {
	Target TargetKey
}

func (s *SingleTarget) IntersectCircle(in Circle) bool       { return false }
func (s *SingleTarget) IntersectRectangle(in Rectangle) bool { return false }
func (s *SingleTarget) Pos() (float64, float64)              { return 0, 0 }
func (s *SingleTarget) String() string                       { return fmt.Sprintf("single target: %v", s.Target) }
