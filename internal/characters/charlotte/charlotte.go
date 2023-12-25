package charlotte

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Charlotte, NewChar)
}

type char struct {
	*tmpl.Character
	markCount int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// TODO subscribe OnTargetDied clear mark
	c.a1()
	c.a4()

	if c.Base.Cons >= 1 {
		c.c1()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Cons >= 6 {
		c.c6()
	}

	return nil
}
