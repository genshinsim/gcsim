package lauma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

type paleHymnStack struct {
	endFrame int
	isC6     bool
}

func init() {
	burstFrames = frames.InitAbilSlice(115) // Q -> walk
	burstFrames[action.ActionAttack] = 111
	burstFrames[action.ActionCharge] = 110
	burstFrames[action.ActionSkill] = 110
	burstFrames[action.ActionDash] = 111
	burstFrames[action.ActionWalk] = 111
	burstFrames[action.ActionSwap] = 109
}

const (
	burstKey          = "lauma-burst"
	burstDeerFrame    = 115
	paleHymnGainFrame = 96
)

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.AddStatus(c1Key, 20*60, true)

	c.Core.Tasks.Add(func() {
		c.addPaleHymn(18, false)
		c.AddStatus(burstKey, 15*60, true) // should this be here?
	}, paleHymnGainFrame)

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{
		Frames: func(next action.Action) int {
			if c.deerStateReady && next == action.ActionCharge {
				return burstDeerFrame
			}
			return burstFrames[next]
		},
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) setupPaleHymnBuff() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if len(c.paleHymnStacks) == 0 {
			return false
		}

		ae := args[1].(*info.AttackEvent)
		em := c.Stat(attributes.EM)

		switch ae.Info.AttackTag {
		case attacks.AttackTagBurningDamage | attacks.AttackTagBloom | attacks.AttackTagHyperbloom | attacks.AttackTagBurgeon:
			ae.Info.FlatDmg += em * bloomDmgIncrease[c.TalentLvlBurst()]

			if c.Base.Cons >= 2 {
				ae.Info.FlatDmg += em * 5
			}
		case attacks.AttackTagDirectLunarBloom:

			// check if lauma c6 skill gets buffed despite not consuming hymn

			ae.Info.FlatDmg += em * lunarBloomDmgIncrease[c.TalentLvlBurst()]

			if c.Base.Cons >= 2 {
				ae.Info.FlatDmg += em * 4
			}
		default:
			return false
		}

		if ae.Info.Abil == "Frostgrove Sanctuary C6" {
			return false
		}

		c.paleHymnStacks = c.paleHymnStacks[:1]

		return false
	}, "lauma-pale-hymn-buff")
}

func (c *char) c6PaleHymn(a info.AttackCB) {
	if c.Base.Cons < 6 {
		return
	}

	c.addPaleHymn(2, true)
}

func (c *char) addPaleHymn(amount int, isC6 bool) {
	if len(c.paleHymnStacks) == 0 {
		c.Core.Tasks.Add(c.removePaleHymn(), 15*60)
	}

	for range amount {
		phs := paleHymnStack{}
		phs.endFrame = c.Core.F + 15*60
		phs.isC6 = false
		c.paleHymnStacks = append(c.paleHymnStacks, phs)
	}

	if isC6 {
		var tmpPaleHymnStacks []paleHymnStack
		for _, phs := range c.paleHymnStacks {
			if phs.isC6 {
				phs.endFrame = c.Core.F + 15*60
			}
			tmpPaleHymnStacks = append(tmpPaleHymnStacks, phs)
		}
		c.paleHymnStacks = tmpPaleHymnStacks
	}
}

func (c *char) removePaleHymn() func() {
	return func() {
		currentFrame := c.Core.F
		var tmpPaleHymnStacks []paleHymnStack
		nextRemovePaleHymn := 0
		for _, phs := range c.paleHymnStacks {
			if currentFrame < phs.endFrame {
				tmpPaleHymnStacks = append(tmpPaleHymnStacks, phs)
			} else {
				if phs.endFrame-currentFrame < nextRemovePaleHymn || nextRemovePaleHymn == 0 {
					nextRemovePaleHymn = phs.endFrame - currentFrame
				}
			}
		}

		if nextRemovePaleHymn != 0 {
			c.Core.Tasks.Add(c.removePaleHymn(), nextRemovePaleHymn)
		}

		c.paleHymnStacks = tmpPaleHymnStacks
	}
}

func (c *char) removeC6PaleHymn() func() {
	return func() {
		var tmpPaleHymnStacks []paleHymnStack

		for _, phs := range c.paleHymnStacks {
			if !phs.isC6 {
				tmpPaleHymnStacks = append(tmpPaleHymnStacks, phs)
			}
		}

		c.paleHymnStacks = tmpPaleHymnStacks
	}
}
