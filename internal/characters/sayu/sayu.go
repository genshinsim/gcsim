package sayu

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl

	eInfused            core.EleType
	eDuration           int
	infuseCheckLocation core.AttackPattern
	c2Bonus             float64
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
	c.c2Bonus = .0

	c.absorbCheck()
	c.a1()

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return &c, nil
}

func (c *char) absorbCheck() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
			return false
		}
		if atk.Info.Element != core.Anemo {
			return false
		}
		if c.Core.F > c.eDuration {
			return false
		}
		if c.eInfused == core.NoElement {
			// TODO: need to check yourself element first
			c.eInfused = c.Core.AbsorbCheck(c.infuseCheckLocation, core.Pyro, core.Hydro, core.Electro, core.Cryo)
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
		case core.AttackTagElementalArt:
			ai := core.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Yoohoo Art: Fuuin Dash (Elemental DMG)",
				AttackTag:  core.AttackTagElementalArt,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    c.eInfused,
				Durability: 25,
				Mult:       skillAbsorb[c.TalentLvlSkill()],
			}
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 1, 1)
		case core.AttackTagElementalArtHold:
			ai := core.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Yoohoo Art: Fuuin Dash (Elemental DMG)",
				AttackTag:  core.AttackTagElementalArt,
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

func (c *char) a1() {
	swirlfunc := func(ele core.EleType) func(args ...interface{}) bool {
		icd := -1
		return func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			if c.Core.ActiveChar != c.Index {
				return false
			}
			if c.Core.F < icd {
				return false
			}
			icd = c.Core.F + 120 // 2s

			if c.Base.Cons >= 4 {
				c.AddEnergy("sayu-c4", 1.2)
			}

			heal := 300 + c.Stat(core.EM)*1.2
			c.Core.Health.Heal(core.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Someone More Capable",
				Src:     heal,
				Bonus:   c.Stat(core.Heal),
			})

			return false
		}
	}

	c.Core.Events.Subscribe(core.OnSwirlCryo, swirlfunc(core.Cryo), "sayu-a1-cryo")
	c.Core.Events.Subscribe(core.OnSwirlElectro, swirlfunc(core.Electro), "sayu-a1-electro")
	c.Core.Events.Subscribe(core.OnSwirlHydro, swirlfunc(core.Hydro), "sayu-a1-hydro")
	c.Core.Events.Subscribe(core.OnSwirlPyro, swirlfunc(core.Pyro), "sayu-a1-pyro")
}

func (c *char) c2() {
	m := make([]float64, core.EndStatType)
	c.AddPreDamageMod(core.PreDamageMod{
		Key: "sayu-c2",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.ActorIndex != c.Index {
				return nil, false
			}
			if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
				return nil, false
			}
			m[core.DmgP] = c.c2Bonus
			c.c2Bonus = .0
			return m, true
		},
		Expiry: -1,
	})
}
