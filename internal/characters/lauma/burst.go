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
	c.DeleteStatus(burstKey)
	c.DeleteStatus(moonSongIcdKey)
	c.Core.Tasks.Add(func() {
		c.c1OnBurst()
		c.addPaleHymn(18)

		c.AddStatus(burstKey, 15*60, true) // should this be here?
		c.moonSongOnBurst()
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
		if c.paleHymnCount() <= 0 {
			return false
		}
		c.consumePaleHymn()

		ae := args[1].(*info.AttackEvent)
		em := c.Stat(attributes.EM)

		switch ae.Info.AttackTag {
		case attacks.AttackTagBountifulCore | attacks.AttackTagBloom | attacks.AttackTagHyperbloom | attacks.AttackTagBurgeon:
			ae.Info.FlatDmg += em * bloomDmgIncrease[c.TalentLvlBurst()]
			ae.Info.FlatDmg += em * c.c2PaleHymnScalingNonLunar()
		case attacks.AttackTagDirectLunarBloom:
			ae.Info.FlatDmg += em * lunarBloomDmgIncrease[c.TalentLvlBurst()]
			ae.Info.FlatDmg += em * c.c2PaleHymnScalingLunar()
		default:
			return false
		}

		if ae.Info.Abil == c6SkillHitName {
			return false
		}

		return false
	}, "lauma-pale-hymn-buff")
}

func (c *char) paleHymnCount() int {
	c.removeExpiredPaleHymn()
	return c.paleHymnStacks.Len() + c.c6PaleHymnCount
}

func (c *char) addPaleHymn(amount int) {
	endFrame := c.TimePassed + 15*60

	for range amount {
		c.paleHymnStacks.PushBack(endFrame)
	}
	c.Core.Tasks.Add(c.removeExpiredPaleHymn, 15*60+1)
}

func (c *char) addC6PaleHymn(amount int) {
	c.c6PaleHymnCount += amount
	c.c6PaleHymnExpiry = c.TimePassed + 15*60
	c.Core.Tasks.Add(c.removeExpiredPaleHymn, 15*60+1)
}

// attempts to consume a pale hymn.
func (c *char) consumePaleHymn() {
	c.removeExpiredPaleHymn()
	if c.paleHymnStacks.Len() == 0 && c.c6PaleHymnCount == 0 {
		// error?
		// panic("consumePaleHymn() called on without Pale Hymn stacks")
		return
	}

	if c.paleHymnStacks.Len() == 0 {
		c.c6PaleHymnCount--
		return
	}
	if c.c6PaleHymnCount <= 0 {
		c.paleHymnStacks.PopFront()
		return
	}

	// pop whichever one is closer to expiry
	currentPaleHymn := c.paleHymnStacks.Front()
	currentC6PaleHymn := c.c6PaleHymnExpiry

	if currentPaleHymn < currentC6PaleHymn {
		c.paleHymnStacks.PopFront()
	} else {
		c.c6PaleHymnCount--
	}
}

func (c *char) removeExpiredPaleHymn() {
	currentFrame := c.TimePassed

	for c.paleHymnStacks.Len() > 0 {
		if c.paleHymnStacks.Front() < currentFrame {
			c.paleHymnStacks.PopFront()
		}
	}

	if c.c6PaleHymnExpiry < currentFrame {
		c.c6PaleHymnCount = 0
	}
}
