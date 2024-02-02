package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c6ICDKey = "wanderer-c6-icd"
)

func (c *char) c1() {
	// C1: Needs to be manually deleted when Windfavored state ends
	if c.Base.Cons < 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.1
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag("wanderer-c1-atkspd", 1200),
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) c2() {
	// C2: Buff stays active during entire burst animation
	if c.Base.Cons < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = min(float64(c.maxSkydwellerPoints-c.skydwellerPoints)*0.04, 2)
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag("wanderer-c2-burstbonus", burstFramesE[action.InvalidAction]),
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) makeC6Callback() func(cb combat.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}

	done := false

	return func(a combat.AttackCB) {
		if done || !c.StatusIsActive(SkillKey) || c.skydwellerPoints <= 0 {
			return
		}

		done = true

		if c.c6Count < 5 && !c.StatusIsActive(c6ICDKey) && c.skydwellerPoints < 40 {
			c.AddStatus(c6ICDKey, 12, true)
			c.c6Count++
			c.skydwellerPoints += 4

			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"wanderer c6 added 4 skydweller points",
			)
		}

		// a gets passed into the callback as param by core
		trg := a.Target

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Shugen: The Curtains’ Melancholic Sway",
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagWandererC6,
			ICDGroup:   attacks.ICDGroupWandererC6,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       a.AttackEvent.Info.Mult * 0.4,
		}

		// TODO: Snapshot delay?
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(trg, nil, 2), 8, 8,
		)
	}
}
