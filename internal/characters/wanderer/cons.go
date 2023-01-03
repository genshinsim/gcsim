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
		c.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("wanderer-c1-atkspd", 1200),
			Amount: func() ([]float64, bool) {
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
		c.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("wanderer-c2-burstbonus", burstFramesE[action.InvalidAction]),
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

	}
}

func (c *char) c6() {
	if c.Base.Cons >= 6 {
		c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
			ae := args[1].(*combat.AttackEvent)
			if ae.Info.ActorIndex != c.Index || ae.Info.Abil == "Shugen: The Curtains’ Melancholic Sway" {
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
					"wanderer c6 added 4 skydweller points",
				)
			}

			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Shugen: The Curtains’ Melancholic Sway",
				AttackTag:  combat.AttackTagNormal,
				ICDTag:     combat.ICDTagWandererC6,
				ICDGroup:   combat.ICDGroupWandererC6,
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
}
