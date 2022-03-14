package kokomi

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Kokomi, NewChar)
}

type char struct {
	*character.Tmpl
	skillLastUsed int
	swapEarlyF    int
	c4ICDExpiry   int
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 70
	}
	c.Energy = float64(e)
	c.EnergyMax = 70
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 3
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneInazuma

	c.skillLastUsed = 0
	c.swapEarlyF = 0
	c.c4ICDExpiry = 0

	c.passive()
	c.onExitField()
	c.burstActiveHook()
	c.a4()

	return &c, nil
}

// Passive 2 - permanently modify stats for +25% healing bonus and -100% CR
func (c *char) passive() {
	val := make([]float64, core.EndStatType)
	val[core.Heal] = .25
	val[core.CR] = -1
	c.AddMod(coretype.CharStatMod{
		Key:    "kokomi-passive",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})
}

func (c *char) a4() {
	c.Core.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if c.Core.StatusDuration("kokomiburst") == 0 {
			return false
		}

		a4Bonus := c.Stat(core.Heal) * 0.15 * c.HPMax
		atk.Info.FlatDmg += a4Bonus

		return false
	}, "kokomi-a4")
}

// Implements event handler for healing during burst
// Also checks constellations
func (c *char) burstActiveHook() {
	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if c.Core.StatusDuration("kokomiburst") == 0 {
			return false
		}

		switch atk.Info.AttackTag {
		case coretype.AttackTagNormal, coretype.AttackTagExtra:
		default:
			return false
		}

		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Ceremonial Garment",
			Src:     burstHealPct[c.TalentLvlBurst()]*c.HPMax + burstHealFlat[c.TalentLvlBurst()],
			Bonus:   c.Stat(core.Heal),
		})

		// C2 handling
		// Sangonomiya Kokomi gains the following Healing Bonuses with regard to characters with 50% or less HP via the following methods:
		// Nereid's Ascension Normal and Charged Attacks: 0.6% of Kokomi's Max HP.
		if c.Base.Cons >= 2 {
			for i, char := range c.Core.Chars {
				if char.HP()/char.MaxHP() > .5 {
					continue
				}
				c.Core.Health.Heal(core.HealInfo{
					Caller:  c.Index,
					Target:  i,
					Message: "The Clouds Like Waves Rippling",
					Src:     0.006 * c.HPMax,
					Bonus:   c.Stat(core.Heal),
				})
			}
		}

		// C4 (Energy piece only) handling
		// While donning the Ceremonial Garment created by Nereid's Ascension, Sangonomiya Kokomi's Normal Attack SPD is increased by 10%, and Normal Attacks that hit opponents will restore 0.8 Energy for her. This effect can occur once every 0.2s.
		if c.Base.Cons >= 4 {
			if c.c4ICDExpiry <= c.Core.Frame {
				c.AddEnergy("kokomi-c4", 0.8)
				c.c4ICDExpiry = c.Core.Frame + 12
			}
		}

		// C6 handling
		// While donning the Ceremonial Garment created by Nereid's Ascension, Sangonomiya Kokomi gains a 40% Hydro DMG Bonus for 4s after her Normal and Charged Attacks heal a character with 80% or more HP.
		if c.Base.Cons == 6 {
			for _, char := range c.Core.Chars {
				if char.HP()/char.MaxHP() < .8 {
					continue
				}
				val := make([]float64, core.EndStatType)
				val[core.HydroP] = .4
				c.AddMod(coretype.CharStatMod{
					Key: "kokomi-c6",
					Amount: func() ([]float64, bool) {
						return val, true
					},
					Expiry: c.Core.Frame + 480,
				})
				// No need to continue checking if we found one
				break
			}
		}

		return false
	}, "kokomi-q-healing")
}

// Clears Kokomi burst when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev == c.Index {
			c.swapEarlyF = c.Core.F
		}
		c.Core.Status.DeleteStatus("kokomiburst")
		return false
	}, "kokomi-exit")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionCharge:
		return 50
	case core.ActionDash:
		return 18
	default:
		c.coretype.Log.NewEvent("ActionStam not implemented", coretype.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
