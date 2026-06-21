package columbina

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	lunarBonusKey        = "columbina-lunar-bonus"
	a1Key                = "columbina-a1"
	moonridgeDewICDKey   = "moonridge-dew-icd"
	moonridgeDewTimerKey = "moonridge-dew-timer"
)

func (c *char) moonsignInit() {
	c.Core.Flags.Custom[reactable.LunarChargeEnableKey] = 1
	c.Core.Flags.Custom[reactable.LunarBloomEnableKey] = 1
	c.Core.Flags.Custom[reactable.LunarCrystallizeEnableKey] = 1
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if !attacks.AttackTagIsLunar(atk.Info.AttackTag) {
			return
		}

		bonus := min(c.MaxHP()/1000.0*0.002, 0.07)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("columbina adding lunar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, lunarBonusKey)
}

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1Buff = make([]float64, attributes.EndStatType)
}

func (c *char) a1OGravityTick() {
	if c.Base.Ascension < 1 {
		return
	}
	if !c.StatModIsActive(a1Key) {
		c.a1Stacks = 1
	} else {
		c.a1Stacks = min(c.a1Stacks+1, 3)
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a1Key, 10*60),
		AffectedStat: attributes.CR,
		Amount: func() []float64 {
			c.a1Buff[attributes.CR] = 0.05 * float64(c.a1Stacks)
			return c.a1Buff
		},
	})
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnLunarCharged, c.a4OnLunarCharge, "columbina-a4-lc")
	c.Core.Events.Subscribe(event.OnLunarBloom, c.a4OnLunarBloom, "columbina-a4-lb")
	c.Core.Events.Subscribe(event.OnLunarCrystallize, c.a4OnLunarCrystallize, "columbina-a4-lcr")
}

func (c *char) a4OnLunarCharge(args ...any) {
	if _, ok := args[0].(*enemy.Enemy); !ok {
		return
	}

	if c.StatusIsActive(burstBuffKey) {
		c.Core.Flags.Custom[reactable.LcIcdOverrideKey] = 1.5 * 60
		return
	}

	// player is outside of lunar domain, reset buffs
	delete(c.Core.Flags.Custom, reactable.LcIcdOverrideKey)
}

func (c *char) a4OnLunarCrystallize(args ...any) {
	if _, ok := args[0].(*enemy.Enemy); !ok {
		return
	}

	if c.StatusIsActive(burstBuffKey) {
		c.Core.Flags.Custom[reactable.LcrExtraHitOverride] = 0.33
		return
	}

	// player is outside of lunar domain, reset buffs
	delete(c.Core.Flags.Custom, reactable.LcrExtraHitOverride)
}

func (c *char) a4OnLunarBloom(args ...any) {
	if _, ok := args[0].(*enemy.Enemy); !ok {
		return
	}

	if !c.StatusIsActive(burstBuffKey) {
		return
	}

	if c.StatusIsActive(moonridgeDewICDKey) {
		return
	}

	if c.a4MoondewCount < 3 {
		c.AddStatus(moonridgeDewICDKey, 0.05*60, true)
		c.a4MoondewCount += 1
		c.Core.Player.AddMoonridgeDew()

		if !c.StatusIsActive(moonridgeDewTimerKey) {
			c.AddStatus(moonridgeDewTimerKey, 18*60, true)
			c.QueueCharTask(func() {
				c.a4MoondewCount = 0
				c.Core.Combat.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), "moonridge dew timer reset")
			}, 18*60)
		}
	}

	// TODO: Moonridge Dew are removed after 60s of not adding any
}
