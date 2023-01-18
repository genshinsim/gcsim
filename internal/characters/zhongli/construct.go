package zhongli

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
)

type stoneStele struct {
	src    int
	expiry int
	c      *char
	dir    combat.Point
	pos    combat.Point
}

func (s *stoneStele) OnDestruct() {
	if s.c.steleCount > 0 {
		s.c.steleCount--
	}
}

func (s *stoneStele) Key() int                         { return s.src }
func (s *stoneStele) Type() construct.GeoConstructType { return construct.GeoConstructZhongliSkill }
func (s *stoneStele) Expiry() int                      { return s.expiry }
func (s *stoneStele) IsLimited() bool                  { return true }
func (s *stoneStele) Count() int                       { return 1 }
func (s *stoneStele) Direction() combat.Point          { return s.dir }
func (s *stoneStele) Pos() combat.Point                { return s.pos }
