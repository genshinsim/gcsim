package zhongli

import "github.com/genshinsim/gcsim/pkg/core"

type stoneStele struct {
	src    int
	expiry int
	c      *char
}

func (s *stoneStele) Key() int {
	return s.src
}

func (s *stoneStele) Type() core.GeoConstructType {
	return core.GeoConstructZhongliSkill
}

func (s *stoneStele) OnDestruct() {
	if s.c.steleCount > 0 {
		s.c.steleCount--
	}
}

func (s *stoneStele) Expiry() int {
	return s.expiry
}

func (s *stoneStele) IsLimited() bool {
	return true
}

func (s *stoneStele) Count() int {
	return 1
}
