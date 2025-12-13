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

var paleHymnStacks = 0

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(111) // Q -> walk
	burstFrames[action.ActionCharge] = 110
	burstFrames[action.ActionSkill] = 110
	burstFrames[action.ActionJump] = 112
	burstFrames[action.ActionSwap] = 109
}

const (
	burstKey = "lauma-burst"
)

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.AddStatus(c1Key, 20*60, true)

	c.Core.Tasks.Add(func() {
		paleHymnStacks += 18
		c.AddStatus(burstKey, 15*60, true)

		c.Core.Tasks.Add(func() {
			paleHymnStacks = 0
		}, 15*60)
	}, 96)

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
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
		if paleHymnStacks <= 0 {
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
			if ae.Info.Abil == "Frostgrove Sanctuary C6" {
				return false
			}

			ae.Info.FlatDmg += em * lunarBloomDmgIncrease[c.TalentLvlBurst()]

			if c.Base.Cons >= 2 {
				ae.Info.FlatDmg += em * 4
			}
		default:
			return false
		}

		paleHymnStacks--

		return false
	}, "lauma-pale-hymn-buff")
}

func (c *char) c6PaleHymn(a info.AttackCB) {
	if c.Base.Cons < 6 {
		return
	}
	paleHymnStacks += 2
}
