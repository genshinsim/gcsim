package kokomi

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Kokomi, NewChar)
}

type char struct {
	*character.Tmpl
	skillLastUsed int
	skillLastTick int
	c4ICDExpiry   int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro
	c.Energy = 70
	c.EnergyMax = 70
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 3
	c.BurstCon = 3
	c.SkillCon = 5
	c.c4ICDExpiry = 0
	c.CharZone = core.ZoneInazuma

	c.passive()
	c.onExitField()
	c.burstActiveHook()

	return &c, nil
}

// Passive 2 - permanently modify stats for +25% healing bonus and -100% CR
func (c *char) passive() {
	val := make([]float64, core.EndStatType)
	val[core.Heal] = .25
	val[core.CR] = -1
	c.AddMod(core.CharStatMod{
		Key:    "kokomi-passive",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, true
		},
	})
}

// Implements event handler for healing during burst
// Also checks constellations
func (c *char) burstActiveHook() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if c.Core.Status.Duration("kokomiburst") == 0 {
			return false
		}

		switch atk.Info.AttackTag {
		case core.AttackTagNormal, core.AttackTagExtra:
		default:
			return false
		}

		c.Core.Health.HealAll(c.Index, burstHealPct[c.TalentLvlBurst()]*c.HPMax+burstHealFlat[c.TalentLvlBurst()])

		// C2 handling
		// Sangonomiya Kokomi gains the following Healing Bonuses with regard to characters with 50% or less HP via the following methods:
		// Nereid's Ascension Normal and Charged Attacks: 0.6% of Kokomi's Max HP.
		if c.Base.Cons >= 2 {
			for i, char := range c.Core.Chars {
				if char.HP()/char.MaxHP() > .5 {
					continue
				}
				c.Core.Health.HealIndex(c.Index, i, 0.006*c.HPMax)
			}
		}

		// C4 (Energy piece only) handling
		// While donning the Ceremonial Garment created by Nereid's Ascension, Sangonomiya Kokomi's Normal Attack SPD is increased by 10%, and Normal Attacks that hit opponents will restore 0.8 Energy for her. This effect can occur once every 0.2s.
		if c.Base.Cons >= 4 {
			if c.c4ICDExpiry <= c.Core.F {
				c.AddEnergy(0.8)
				c.c4ICDExpiry = c.Core.F + 12
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
				c.AddMod(core.CharStatMod{
					Key: "kokomi-c6",
					Amount: func(a core.AttackTag) ([]float64, bool) {
						return val, true
					},
					Expiry: c.Core.F + 480,
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
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Key.String())
		return 0
	}
}
