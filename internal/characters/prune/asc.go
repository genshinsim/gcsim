package prune

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1SwirlCheckICDKey = "prune-a1-swirl-check-icd"
	a4Key              = "prune-tolling-rally"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	swirlfunc := func(ele attributes.Element) func(args ...any) {
		return func(args ...any) {
			// reaction target must be an enemy
			target, ok := args[0].(*enemy.Enemy)
			if !ok {
				return
			}

			// reaction must have been triggered by prune's attack
			atk, ok := args[1].(*info.AttackEvent)
			if !ok {
				return
			}

			if atk.Info.ActorIndex != c.Index() {
				return
			}

			// reaction must have been triggered by the burst dot
			if atk.Info.Abil != "Witchlure Bell" {
				return
			}

			// A1 has a shared 1.2s ICD across all four swirl elements
			if c.StatusIsActive(a1SwirlCheckICDKey) {
				return
			}
			c.AddStatus(a1SwirlCheckICDKey, 72, false)

			// queue passive
			ai := info.AttackInfo{
				ActorIndex: c.Index(),
				Abil:       "Verdict and Punishment (A1)",
				AttackTag:  attacks.AttackTagElementalBurst,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeBlunt,
				Element:    ele,
				Durability: 0,
				Mult:       1.5,
			}

			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(target, nil, 2.3),
				0,
				45,
				c.makeA4CB,
				c.makeC1CB,
				c.c4Ricochet(c.a1ConvertEle, attacks.AttackTagElementalBurst),
				c.makeC2CB,
			)

			// swirl burst has different energy drain delay frame
			c.burstEnergyDrainDelay = 19

		}
	}

	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc(attributes.Cryo), "prune-burst-cryo")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc(attributes.Electro), "prune-burst-electro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc(attributes.Hydro), "prune-burst-hydro")
	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc(attributes.Pyro), "prune-burst-pyro")

}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.a4Buff = make([]float64, attributes.EndStatType)
	atk := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK).TotalATK()
	c.a4Buff[attributes.DmgP] = math.Min(math.Max(atk-2000, 0)*0.00025, 0.5)

	for i, char := range c.Core.Player.Chars() {
		if i == c.Index() {
			continue // nothing for prune, need testing
		}
		char.AddStatus(a4Key, 5*60, true)
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("prune-a4", 5*60), // 5 s
			Amount: func(
				atk *info.AttackEvent,
				target info.Target,
			) []float64 {
				return c.a4Buff
			},
		})

	}

	c.Core.Log.NewEvent("prune a4 triggered", glog.LogCharacterEvent, c.Index()).
		Write("dmg bonus", c.a4Buff[attributes.DmgP]).
		Write("expiry", c.Core.F+5*60)
}

func (c *char) makeA4CB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	c.a4()
}

func (c *char) hexInit() {
	if !c.IsHexerei {
		return
	}
	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	chars := c.Core.Player.Chars()

	// separate values are needed because a swirl should be able to apply both buffs on the same frame
	lastSelfTriggerFrame := make([]int, len(chars))
	lastTeamTriggerFrame := make([]int, len(chars))

	for i := range chars {
		lastSelfTriggerFrame[i] = -1
		lastTeamTriggerFrame[i] = -1
	}

	selfBuff := make([]float64, attributes.EndStatType)
	selfBuff[attributes.ATKP] = 0.60

	teamBuff := make([]float64, attributes.EndStatType)
	teamBuff[attributes.ATKP] = 0.30

	reactionCB := func(isSwirl bool) func(args ...any) {
		return func(args ...any) {
			if len(args) < 2 {
				return
			}

			// reaction target must be an enemy
			if _, ok := args[0].(*enemy.Enemy); !ok {
				return
			}

			// reaction must have been triggered by a character attack
			atk, ok := args[1].(*info.AttackEvent)
			if !ok {
				return
			}

			triggererIndex := atk.Info.ActorIndex
			if triggererIndex < 0 || triggererIndex >= len(chars) {
				return
			}

			triggerer := chars[triggererIndex]

			// triggerer must be a hexerei character affected by tolling rally
			if !triggerer.IsHexerei || !triggerer.StatusIsActive(a4Key) {
				return
			}

			// any qualifying reaction grants prune's self buff
			if lastSelfTriggerFrame[triggererIndex] != c.Core.F {
				lastSelfTriggerFrame[triggererIndex] = c.Core.F

				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("prune-hex-self-buff", 5*60),
					AffectedStat: attributes.ATKP,
					Amount: func() []float64 {
						return selfBuff
					},
				})

				c.Core.Log.NewEvent("prune hex self buff triggered", glog.LogCharacterEvent, c.Index()).
					Write("triggerer", triggererIndex).
					Write("atk percent", selfBuff[attributes.ATKP]).
					Write("expiry", c.Core.F+5*60)
			}

			// only swirl grants the triggering character's team buff
			if !isSwirl {
				return
			}

			if lastTeamTriggerFrame[triggererIndex] == c.Core.F {
				return
			}
			lastTeamTriggerFrame[triggererIndex] = c.Core.F

			triggerer.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("prune-hex-team-buff", 5*60),
				AffectedStat: attributes.ATKP,
				Amount: func() []float64 {
					return teamBuff
				},
			})

			c.Core.Log.NewEvent("prune hex team buff triggered", glog.LogCharacterEvent, c.Index()).
				Write("triggerer", triggererIndex).
				Write("atk percent", teamBuff[attributes.ATKP]).
				Write("expiry", c.Core.F+5*60)
		}
	}

	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		isSwirl := i == event.OnSwirlCryo || i == event.OnSwirlElectro || i == event.OnSwirlHydro || i == event.OnSwirlPyro
		c.Core.Events.Subscribe(i, reactionCB(isSwirl), "prune-hex")
	}
}
