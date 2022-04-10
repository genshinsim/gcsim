package sayu

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl

	eInfused  core.EleType
	eDuration int
}

func init() {
	core.RegisterCharFunc(core.Sayu, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5

	c.eInfused = core.NoElement
	c.eDuration = -1

	c.absorbCheck()

	return &c, nil
}

func (c *char) absorbCheck() {
	icd := -1
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagSayuRoll && atk.Info.AttackTag != core.AttackTagElementalArtHold {
			return false
		}
		if c.Core.F > c.eDuration {
			return false
		}
		if c.eInfused == core.NoElement {
			// TODO: need to check yourself element first
			c.eInfused = c.Core.AbsorbCheck(core.Pyro, core.Hydro, core.Electro, core.Cryo)
			if c.eInfused == core.NoElement {
				return false
			}

			c.Core.Log.NewEventBuildMsg(
				core.LogCharacterEvent,
				c.Index,
				"sayu infused ", c.eInfused.String(),
			)
		}

		switch atk.Info.AttackTag {
		case core.AttackTagSayuRoll:
			if icd > c.Core.F {
				return false
			}
			ai := core.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Yoohoo Art: Fuuin Dash (Elemental DMG)",
				AttackTag:  core.AttackTagSayuRoll,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    c.eInfused,
				Durability: 25,
				Mult:       skillAbsorb[c.TalentLvlSkill()],
			}
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 1, 1)
			icd = c.Core.F + 30
		case core.AttackTagElementalArtHold:
			ai := core.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Yoohoo Art: Fuuin Dash (Elemental DMG)",
				AttackTag:  core.AttackTagSayuRoll,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    c.eInfused,
				Durability: 25,
				Mult:       skillAbsorbEnd[c.TalentLvlSkill()],
			}
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 1, 1)
		}

		return false
	}, "sayu-absorb-check")
}
