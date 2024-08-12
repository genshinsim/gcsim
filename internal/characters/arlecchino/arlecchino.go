package arlecchino

import (
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Arlecchino, NewChar)
}

type char struct {
	*tmpl.Character
	skillDebt             float64
	skillDebtMax          float64
	initialDirectiveLevel int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = base.SkillDetails.BurstEnergyCost
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.naBuff()
	c.passive()

	c.a1OnKill()
	c.a4()

	c.c2()
	return nil
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	lastAction := c.Character.Core.Player.LastAction
	if k != c.Base.Key && a != action.ActionSwap {
		return fmt.Errorf("%v: Tried to execute %v when not on field", c.Base.Key, a)
	}

	if lastAction.Type == action.ActionCharge && lastAction.Param["early_cancel"] > 0 {
		// can only early cancel charged attack with Dash or Jump
		switch a {
		case action.ActionDash, action.ActionJump: // skips the error in default block
		default:
			return fmt.Errorf("%v: Cannot early cancel Charged Attack with %v", c.Base.Key, a)
		}
	}

	// can use charge without attack beforehand unlike most of the other polearm users
	if a == action.ActionCharge {
		return nil
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 15
	case model.AnimationYelanN0StartDelay:
		return 7
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) ReceiveHeal(hi *info.HealInfo, healAmt float64) float64 {
	// ignore all healing except hers
	if hi.Caller == c.Index {
		return c.Character.ReceiveHeal(hi, healAmt)
	}
	return c.GetReceivedHeal(healAmt)
}
