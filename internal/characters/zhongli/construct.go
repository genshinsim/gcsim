package zhongli

import (
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type stoneStele struct {
	src    int
	expiry int
	c      *char
	dir    info.Point
	pos    info.Point
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
func (s *stoneStele) Direction() info.Point            { return s.dir }
func (s *stoneStele) Pos() info.Point                  { return s.pos }
