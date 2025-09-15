package chongyun

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Chongyun, NewChar)
}

type char struct {
	*tmpl.Character
	skillArea info.AttackPattern
	fieldSrc  int
	a4Snap    info.Snapshot
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.fieldSrc = -601

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onSwapHook()
	if c.Base.Cons >= 6 && c.Core.Combat.DamageMode {
		c.c6()
	}
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 18
	}
	return c.Character.AnimationStartDelay(k)
}
