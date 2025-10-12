package flins

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Flins, NewChar)
}

type char struct {
	*tmpl.Character
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3
	c.Moonsign = 1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Init()
	c.a4Init()
	c.lunarchargeInit()

	c.c1Init()
	c.c2GleamInit()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 15
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if c.StatusIsActive(skillKey) && a == action.ActionSkill {
		if c.StatusIsActive(spearStormCDKey) {
			return false, action.SkillCD
		}
		return true, action.NoFailure
	}

	if c.StatModIsActive(thunderousSymphonyKey) && a == action.ActionBurst {
		if !c.Core.Flags.IgnoreBurstEnergy && c.Energy < 30 {
			return false, action.InsufficientEnergy
		}
		return true, action.NoFailure
	}

	return c.Character.ActionReady(a, p)
}

func (c *char) NextQueueItemIsValid(_ keys.Char, a action.Action, p map[string]int) error {
	if a == action.ActionCharge {
		switch c.Weapon.Class {
		case info.WeaponClassSword, info.WeaponClassSpear:
			if c.NormalCounter == 0 {
				return player.ErrInvalidChargeAction
			}
		}
	}

	return nil
}

func (c *char) getMoonsignLevel() int {
	count := 0
	for _, c := range c.Core.Player.Chars() {
		count += c.Moonsign
	}
	return count
}
