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

const (
	paleHymnC6       = 0
	paleHymnMoonsong = 1
	paleHymnBurst    = 2

	paleHymnC6Key       = "lauma-pale-hymn-c6"
	paleHymnMoonsongKey = "lauma-pale-hymn-moonsong"
	paleHymnBurstKey    = "lauma-pale-hymn-burst"
)

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
		c.setPaleHymnBurst(18)

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
	return c.paleHymn[paleHymnBurst] + c.paleHymn[paleHymnMoonsong] + c.paleHymn[paleHymnC6]
}

func (c *char) setPaleHymnBurst(amount int) {
	ind := paleHymnBurst
	c.paleHymn[ind] = amount
	c.paleHymnSrc[ind] = c.Core.F
	c.AddStatus(paleHymnBurstKey, paleHymnDur, true)
	c.QueueCharTask(c.removeExpiredPaleHymn(c.paleHymnSrc[ind], ind), paleHymnDur)
}

func (c *char) setPaleHymnMoonsong(amount int) {
	ind := paleHymnMoonsong
	c.paleHymn[ind] = amount
	c.paleHymnSrc[ind] = c.Core.F
	c.AddStatus(paleHymnMoonsongKey, paleHymnDur, true)
	c.QueueCharTask(c.removeExpiredPaleHymn(c.paleHymnSrc[ind], ind), paleHymnDur)
}

func (c *char) addC6PaleHymn(amount int) {
	ind := paleHymnC6
	c.paleHymn[ind] += amount
	c.paleHymnSrc[ind] = c.Core.F
	c.AddStatus(paleHymnC6Key, paleHymnDur, true)
	c.QueueCharTask(c.removeExpiredPaleHymn(c.paleHymnSrc[ind], ind), paleHymnDur)
}

// attempts to consume a pale hymn.
func (c *char) consumePaleHymn() {
	if c.paleHymn[paleHymnC6] > 0 {
		c.paleHymn[paleHymnC6]--
		return
	}

	if c.paleHymn[paleHymnMoonsong] > 0 {
		c.paleHymn[paleHymnMoonsong]--
		return
	}

	if c.paleHymn[paleHymnBurst] > 0 {
		c.paleHymn[paleHymnBurst]--
		return
	}

	// err or panic?
	// panic("consumePaleHymn called when there are no pale hymn stacks")
}

func (c *char) removeExpiredPaleHymn(src, index int) func() {
	return func() {
		if c.paleHymnSrc[index] != src {
			return
		}
		c.paleHymn[index] = 0
	}
}
