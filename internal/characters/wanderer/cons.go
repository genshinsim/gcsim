package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"math"
)

const (
	c6ICDKey = "wanderer-c6-icd"
)

func (c *char) c1() {
	// C1: Needs to be manually deleted when Windfavored state ends
	if c.Base.Cons >= 1 && c.StatusIsActive(skillKey) {
		m := make([]float64, attributes.EndStatType)
		m[attributes.AtkSpd] = 0.1
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("wanderer-c1-atkspd", 1200),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if !(atk.Info.AttackTag == combat.AttackTagNormal || atk.Info.AttackTag == combat.AttackTagExtra) {
					return nil, false
				}
				return m, true
			},
		})

	}
}

func (c *char) c2() {
	// C2: Buff stays active during entire burst animation
	if c.Base.Cons >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = math.Min((float64)(c.maxSkydwellerPoints-c.skydwellerPoints)*0.04, 2)
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("wanderer-c2-burstbonus", burstFramesE[action.InvalidAction]),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if !(atk.Info.AttackTag == combat.AttackTagElementalBurst) {
					return nil, false
				}
				return m, true
			},
		})

	}
}

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {

		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}

		switch ae.Info.AttackTag {
		case combat.AttackTagNormal:
		default:
			return false
		}

		if c.c6Count < 5 && !c.StatusIsActive(c6ICDKey) && c.skydwellerPoints < 40 {
			c.AddStatus(c6ICDKey, 12, true)
			c.c6Count++
			c.skydwellerPoints += 4

			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"wanderer c6 added 4 skydweller points ",
			)
		}

		// TODO: ICD info taken from KQM, seems to just count as a normal attack
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Shugen: The Curtains’ Melancholic Sway",
			AttackTag:  combat.AttackTagNormal,
			ICDTag:     combat.ICDTagNormalAttack,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       ae.Info.Mult * 0.4,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.5), 0, 0,
		)

		return false
	}, "wanderer-c6")
}
