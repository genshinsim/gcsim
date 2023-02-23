package combat

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/targets"
)

type SingleTarget struct {
	Target targets.TargetKey
}

func (s *SingleTarget) PointInShape(p Point) bool            { return true }
func (s *SingleTarget) IntersectCircle(in Circle) bool       { return false }
func (s *SingleTarget) IntersectRectangle(in Rectangle) bool { return false }
func (s *SingleTarget) Pos() Point                           { return Point{X: 0, Y: 0} }
func (s *SingleTarget) String() string                       { return fmt.Sprintf("single target: %v", s.Target) }
