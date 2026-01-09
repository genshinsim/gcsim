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
	c.c1OnBurst()

	c.Core.Tasks.Add(func() {
		c.addPaleHymn(18)

		c.AddStatus(burstKey, 15*60, true) // should this be here?
		if c.moonSong != 0 {
			c.addPaleHymn(6 * c.moonSong)
			c.moonSong = 0
			c.AddStatus(moonSongAddedKey, 15*60, true)
		}
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

func (c *char) initBurst() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if c.paleHymnStacks.Len() == 0 && c.c6PaleHymnStacks.Len() == 0 {
			return false
		}

		ae := args[1].(*info.AttackEvent)
		em := c.Stat(attributes.EM)

		switch ae.Info.AttackTag {
		case attacks.AttackTagBountifulCore | attacks.AttackTagBloom | attacks.AttackTagHyperbloom | attacks.AttackTagBurgeon:
			ae.Info.FlatDmg += em * bloomDmgIncrease[c.TalentLvlBurst()]

			ae.Info.FlatDmg += em * c.c2PaleHymnScaling(false)
		case attacks.AttackTagDirectLunarBloom:

			ae.Info.FlatDmg += em * lunarBloomDmgIncrease[c.TalentLvlBurst()]

			ae.Info.FlatDmg += em * c.c2PaleHymnScaling(true)
		default:
			return false
		}

		if ae.Info.Abil == c6SkillHitName {
			return false
		}

		c.consumePaleHymn()

		return false
	}, "lauma-pale-hymn-buff")
}

func (c *char) addPaleHymn(amount int) {
	if c.paleHymnStacks.Len() == 0 {
		c.Core.Tasks.Add(c.removePaleHymn(), 15*60)
	}

	endFrame := c.Core.F + 15*60

	for range amount {
		c.paleHymnStacks.PushBack(endFrame)
	}
}

func (c *char) addC6PaleHymn(amount int) {
	if c.paleHymnStacks.Len() == 0 {
		c.Core.Tasks.Add(c.removePaleHymn(), 15*60)
	}

	endFrame := c.Core.F + 15*60

	for range amount {
		c.c6PaleHymnStacks.PushBack(endFrame)
	}
}

func (c *char) consumePaleHymn() {
	if c.paleHymnStacks.Len() == 0 {
		c.c6PaleHymnStacks.PopFront()
		return
	}
	if c.c6PaleHymnStacks.Len() == 0 {
		c.paleHymnStacks.PopFront()
		return
	}
	currentPaleHymn := c.paleHymnStacks.Front()
	currentC6PaleHymn := c.c6PaleHymnStacks.Front()

	if currentPaleHymn < currentC6PaleHymn {
		c.paleHymnStacks.PopFront()
	} else {
		c.c6PaleHymnStacks.PopFront()
	}
}

func (c *char) removePaleHymn() func() {
	return func() {
		currentFrame := c.Core.F

		var nextRemovePaleHymn int
		var currentPaleHymn int
		var currentC6PaleHymn int

		if c.paleHymnStacks.Len() != 0 {
			currentPaleHymn = c.paleHymnStacks.Front()

			for currentPaleHymn <= currentFrame {
				if c.paleHymnStacks.Len() == 0 {
					break
				}
				c.paleHymnStacks.PopFront()
			}

			if currentPaleHymn > currentFrame {
				nextRemovePaleHymn = currentPaleHymn
			}
		}

		if c.c6PaleHymnStacks.Len() != 0 {
			currentC6PaleHymn = c.c6PaleHymnStacks.Back()

			if currentC6PaleHymn <= currentFrame {
				c.c6PaleHymnStacks.Clear()
			}
		}

		if currentC6PaleHymn > currentFrame && currentC6PaleHymn < nextRemovePaleHymn {
			nextRemovePaleHymn = currentC6PaleHymn
		}

		if nextRemovePaleHymn != 0 {
			c.Core.Tasks.Add(c.removePaleHymn(), nextRemovePaleHymn-currentFrame)
		}
	}
}
