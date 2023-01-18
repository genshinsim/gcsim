package cyno

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Cyno, NewChar)
}

type char struct {
	*tmpl.Character
	burstExtension int
	burstSrc       int
	lastSkillCast  int
	c4Counter      int
	c6Stacks       int
	a1Extended     bool
	normalBCounter int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.a1Extension()

	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Cons >= 6 {
		c.c6()
	}

	return nil
}

func (c *char) AdvanceNormalIndex() {
	if c.StatusIsActive(BurstKey) {
		c.normalBCounter++
		if c.normalBCounter == burstHitNum {
			c.normalBCounter = 0
		}
		return
	}
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *char) ResetNormalCounter() {
	c.normalBCounter = 0
	c.NormalCounter = 0
}

func (c *char) NextNormalCounter() int {
	if c.StatusIsActive(BurstKey) {
		return c.normalBCounter + 1
	}
	return c.NormalCounter + 1
}
