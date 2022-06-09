package lisa

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Lisa, NewChar)
}

type char struct {
	*character.Tmpl
	c6icd int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()
	c.InitCancelFrames()

	c.skillHoldMult()

	if c.Base.Cons == 6 {
		c.c6()
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 50
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}

}

func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.F < c.c6icd && c.c6icd != 0 {
			return false
		}
		if c.Core.ActiveChar == c.CharIndex() {
			//swapped to lisa

			// Create a "fake attack" to apply conductive stacks to all nearby opponents
			// Needed to ensure hitboxes are properly accounted for
			// Similar to current "Freeze Breaking" solution
			ai := core.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Lisa C6 Conductive Status Application",
				AttackTag:  core.AttackTagNone,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.NoElement,
				DoNotLog:   true,
			}
			cb := func(a core.AttackCB) {
				a.Target.SetTag(conductiveTag, 3)
			}
			// TODO: No idea what the exact radius of this is
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), -1, 0, cb)

			c.c6icd = c.Core.F + 300
		}
		return false
	}, "lisa-c6")
}

func (c *char) skillHoldMult() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if atk.Info.Abil != "Violet Arc (Hold)" {
			return false
		}
		stacks := t.GetTag(conductiveTag)

		atk.Info.Mult = skillHold[stacks][c.TalentLvlSkill()]

		//consume the stacks
		t.SetTag(conductiveTag, 0)

		return false
	}, "lisa-skill-hold-mul")
}
