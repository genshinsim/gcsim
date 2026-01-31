package lauma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const paleHymnDur = 15 * 60

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

		c.AddStatus(burstKey, paleHymnDur, true) // should this be here?
		c.moonSongOnBurst()
	}, paleHymnGainFrame)

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, paleHymnDur)
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

		c.consumePaleHymn()
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("lauma pale hymn consumed", glog.LogCharacterEvent, c.Index()).Write("remaining", c.paleHymnCount())
		}

		return false
	}, "lauma-pale-hymn-buff")
}

func (c *char) paleHymnCount() int {
	c.removeExpiredPaleHymn()
	return c.paleHymnStacksSrc.Len() + c.paleHymnStacksSrcC6.Len()
}

func (c *char) addPaleHymn(amount int) {
	startFrame := c.TimePassed

	for range amount {
		c.paleHymnStacksSrc.PushBack(startFrame)
	}
	c.QueueCharTask(c.removeExpiredPaleHymn, paleHymnDur+1)
}

func (c *char) addC6PaleHymn(amount int) {
	startFrame := c.TimePassed
	for range amount {
		c.paleHymnStacksSrcC6.PushBack(startFrame)
	}
	c.c6PaleHymnExpiry = c.TimePassed + paleHymnDur
	c.QueueCharTask(c.removeExpiredPaleHymn, paleHymnDur+1)
}

// attempts to consume a pale hymn.
func (c *char) consumePaleHymn() {
	c.removeExpiredPaleHymn()
	if c.paleHymnStacksSrc.Len() == 0 && c.paleHymnStacksSrcC6.Len() == 0 {
		// error?
		// panic("consumePaleHymn() called on without Pale Hymn stacks")
		return
	}

	if c.paleHymnStacksSrc.Len() == 0 {
		c.paleHymnStacksSrcC6.PopFront()
		return
	}

	if c.paleHymnStacksSrcC6.Len() == 0 {
		c.paleHymnStacksSrc.PopFront()
		return
	}

	// pop whichever one is closer to expiry
	paleHymnSrc := c.paleHymnStacksSrc.Front()
	paleHymnSrcC6 := c.paleHymnStacksSrcC6.Front()

	if paleHymnSrc < paleHymnSrcC6 {
		c.paleHymnStacksSrc.PopFront()
	} else {
		c.paleHymnStacksSrcC6.PopFront()
	}
}

func (c *char) removeExpiredPaleHymn() {
	currentFrame := c.TimePassed

	for c.paleHymnStacksSrc.Len() > 0 && c.paleHymnStacksSrc.Front()+paleHymnDur < currentFrame {
		c.paleHymnStacksSrc.PopFront()
	}

	if c.c6PaleHymnExpiry < currentFrame {
		c.paleHymnStacksSrcC6.Clear()
	}
}
