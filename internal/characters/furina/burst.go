package furina

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark            = 98
	burstDur                = 18.2 * 60
	burstKey                = "furina-burst"
	fanfareDebounceKey      = "furina-fanfare-debounce"
	fanfareDrainToGainDelay = 6                                // estimation
	fanfareDebounceDur      = 30 + fanfareDrainToGainDelay - 1 // estimation
)

func init() {
	burstFrames = frames.InitAbilSlice(121)
	burstFrames[action.ActionAttack] = 113 // Q -> N1
	burstFrames[action.ActionCharge] = 113 // Q -> CA
	burstFrames[action.ActionSkill] = 114  // Q -> E
	burstFrames[action.ActionDash] = 115   // Q -> D
	burstFrames[action.ActionJump] = 115   // Q -> J
	burstFrames[action.ActionSwap] = 111   // Q -> Swap
}

func (c *char) addFanfareFunc(amt float64) func() {
	return func() {
		if c.Base.Cons >= 2 {
			amt *= 3.5
		}
		prevFanfare := c.curFanfare
		c.curFanfare = min(c.maxC2Fanfare, c.curFanfare+amt)
		c.Core.Log.NewEvent("Gained Fanfare", glog.LogCharacterEvent, c.Index).
			Write("previous fanfare", prevFanfare).
			Write("current fanfare", c.curFanfare)
	}
}

func (c *char) queueFanfareGain(amt float64) {
	// determine delay between hp drain and actual fanfare change based on debounce status being up or not
	var delay int
	if !c.StatusIsActive(fanfareDebounceKey) {
		// leave a 1f window open for drains of other chars in same frame to queue fanfare gain until debounce status is added
		// use a char bool to make sure only 1 debounce status add task is queued at a time
		if !c.fanfareDebounceTaskQueued {
			c.fanfareDebounceTaskQueued = true
			c.Core.Tasks.Add(func() {
				c.AddStatus(fanfareDebounceKey, fanfareDebounceDur, false) // TODO: unsure about hitlag
				c.fanfareDebounceTaskQueued = false
			}, 1)
		}
		// fanfare gain from drain has a delay even after drain is confirmed via server
		delay = fanfareDrainToGainDelay
	} else {
		// fanfare gain from drain has to be delayed until debounce status is gone
		delay = c.StatusDuration(fanfareDebounceKey)
	}
	// queue fanfare change
	c.Core.Tasks.Add(c.addFanfareFunc(amt), delay)
}

func (c *char) burstInit() {
	c.maxC2Fanfare = 300
	c.maxQFanfare = 300
	if c.Base.Cons >= 1 {
		c.maxQFanfare = 400
		c.maxC2Fanfare = 400
	}
	if c.Base.Cons >= 2 {
		// 400 + 140/0.35 = 800
		c.maxC2Fanfare = 800
	}
	c.burstBuff = make([]float64, attributes.EndStatType)

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}

		di := args[0].(player.DrainInfo)

		if di.Amount <= 0 {
			return false
		}

		char := c.Core.Player.ByIndex(di.ActorIndex)
		amt := di.Amount / char.MaxHP() * 100
		c.queueFanfareGain(amt)

		return false
	}, "furina-fanfare-on-hp-drain")

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}

		target := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)

		if amount <= 0 {
			return false
		}

		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}

		char := c.Core.Player.ByIndex(target)
		amt := (amount - overheal) / char.MaxHP() * 100

		c.queueFanfareGain(amt)

		return false
	}, "furina-fanfare-on-heal")

	burstDMGRatio := burstFanfareDMGRatio[c.TalentLvlBurst()]
	burstHealRatio := burstFanfareHBRatio[c.TalentLvlBurst()]
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("furina-burst-damage-buff", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if !c.StatusIsActive(burstKey) {
					return nil, false
				}
				c.burstBuff[attributes.DmgP] = min(c.curFanfare, c.maxQFanfare) * burstDMGRatio
				return c.burstBuff, true
			},
		})

		char.AddHealBonusMod(character.HealBonusMod{
			Base: modifier.NewBase("furina-burst-heal-buff", -1),
			Amount: func() (float64, bool) {
				if c.StatusIsActive(burstKey) {
					return min(c.curFanfare, c.maxQFanfare) * burstHealRatio, false
				}
				return 0, false
			},
		})
	}
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Let the People Rejoice",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    c.MaxHP() * burstDMG[c.TalentLvlBurst()],
	}

	c.curFanfare = 0
	c.DeleteStatus(burstKey)

	c.QueueCharTask(func() {
		if c.Base.Cons >= 1 {
			c.curFanfare = 150
		}
		c.AddStatus(burstKey, burstDur, true)
	}, 95) // This is before the hitmark, so Furina burst damage will benefit from C1. This was tested and confirmed

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), burstHitmark, burstHitmark)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
