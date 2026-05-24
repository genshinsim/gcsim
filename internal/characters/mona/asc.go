package mona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
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

var (
	astralGlowICDKey       = "mona-astral-glow-icd"
	omenRefreshICDKey      = "mona-omen-refresh-icd"
	omenRefreshCount       int
	astralGlowExpiryFrames []int
)

// After she has used Illusory Torrent for 2s, if there are any opponents nearby,
// Mona will automatically create a Phantom.
// A Phantom created in this manner lasts for 2s, and its explosion DMG is equal to 50% of Mirror Reflection of Doom.
//
// - checks for ascension level in dash.go to avoid queuing this up only to fail the ascension level check
func (c *char) a1() {
	// do nothing if not Mona
	if c.Core.Player.Active() != c.Index() {
		return
	}
	// do nothing if we aren't dashing anymore
	if c.Core.Player.CurrentState() != action.DashState {
		return
	}
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 15), nil)
	if enemies != nil {
		c.Core.Log.NewEvent("mona-a1 phantom added", glog.LogCharacterEvent, c.Index()).
			Write("expiry:", c.Core.F+120)
		// queue up phantom explosion
		phantomPos := c.Core.Combat.Player()
		c.Core.Tasks.Add(func() {
			aiExplode := info.AttackInfo{
				ActorIndex: c.Index(),
				Abil:       "Mirror Reflection of Doom (A1 Explode)",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Hydro,
				Durability: 25,
				Mult:       0.5 * skill[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(aiExplode, combat.NewCircleHitOnTarget(phantomPos, nil, 5), 0, 0)
		}, 120)
	}
	// queue up next A1 check because Mona's still dashing
	// different Phantoms coexist and don't overwrite each other
	c.Core.Tasks.Add(c.a1, 120) // check again in 2s
}

// Increases Mona's Hydro DMG Bonus by a degree equivalent to 20% of her Energy Recharge rate.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	if c.a4Stats == nil {
		c.a4Stats = make([]float64, attributes.EndStatType)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("mona-a4", -1),
			AffectedStat: attributes.HydroP,
			Extra:        true,
			Amount: func() []float64 {
				return c.a4Stats
			},
		})
	}
	c.a4Stats[attributes.HydroP] = 0.2 * c.NonExtraStat(attributes.ER)
	c.QueueCharTask(c.a4, 60)
}

func (c *char) hexInit() {
	if !c.IsHexerei {
		return
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	astralGlowGainHook := func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)

		if !ok {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if c.StatusIsActive(astralGlowICDKey) {
			return
		}

		c.AddStatus(astralGlowICDKey, 6, false) // 0.1s ICD

		astralGlowExpiryFrames = append(astralGlowExpiryFrames, c.Core.F+60*8) // 8s duration

		if len(astralGlowExpiryFrames) > 3 {
			astralGlowExpiryFrames = astralGlowExpiryFrames[1:]
		}

		c.Core.Log.NewEvent("mona hexerei proc: astral glow", glog.LogCharacterEvent, c.Index()).
			Write("expiry:", c.Core.F+60*8)

		c.Core.Tasks.Add(c.removeAstralGlowStack, 60*8)
	}

	omenRefreshHook := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)

		if !ok {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if omenRefreshCount > 3 {
			return
		}

		if !t.StatusIsActive(omenKey) {
			return
		}

		if c.StatusIsActive(omenRefreshICDKey) {
			return
		}

		c.AddStatus(omenRefreshICDKey, 30, false) // 0.5s ICD

		t.AddStatus(omenKey, t.StatusExpiry(omenKey)+60*2, true)

		omenRefreshCount++

		c.Core.Log.NewEvent("mona hexerei proc: omen refresh", glog.LogCharacterEvent, c.Index()).
			Write("refreshCount", omenRefreshCount)
	}

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("mona-hexerei-astral-glow-vaporize", -1),
			Amount: func(ai info.AttackInfo) float64 {
				m := 0.05 * float64(len(astralGlowExpiryFrames))

				if ai.Amped {
					astralGlowExpiryFrames = []int{} // clear all stacks
					return m
				}

				return 0
			},
		})
	}

	c.AddReactBonusMod(character.ReactBonusMod{})

	c.Core.Events.Subscribe(event.OnEnemyDamage, astralGlowGainHook, "mona-hex-astral-glow")
	c.Core.Events.Subscribe(event.OnEnemyDamage, omenRefreshHook, "mona-hex-omen-refresh")
}

func (c *char) removeAstralGlowStack() {
	if len(astralGlowExpiryFrames) == 0 {
		return
	}

	if c.Core.F >= astralGlowExpiryFrames[0] {
		astralGlowExpiryFrames = astralGlowExpiryFrames[1:]
		c.Core.Log.NewEvent("mona hexerei expired: astral glow", glog.LogCharacterEvent, c.Index())
	}
}
