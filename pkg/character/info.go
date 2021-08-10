package character

import "github.com/genshinsim/gsim/pkg/core"

func (c *Tmpl) Tag(key string) int {
	return c.Tags[key]
}

func (c *Tmpl) CharIndex() int {
	return c.Index
}

func (c *Tmpl) Name() string {
	return c.Base.Name
}

func (c *Tmpl) Zone() core.ZoneType {
	return c.CharZone
}

func (c *Tmpl) Ele() core.EleType {
	return c.Base.Element
}

func (c *Tmpl) WeaponClass() core.WeaponClass {
	return c.Weapon.Class
}

func (c *Tmpl) TalentLvlSkill() int {
	if c.Base.Cons >= c.SkillCon {
		return c.Talents.Skill + 2
	}
	return c.Talents.Skill - 1
}
func (c *Tmpl) TalentLvlBurst() int {
	if c.Base.Cons >= c.BurstCon {
		return c.Talents.Burst + 2
	}
	return c.Talents.Burst - 1
}
func (c *Tmpl) TalentLvlAttack() int {
	if c.Sim.Flags().ChildeActive {
		return c.Talents.Attack
	}
	return c.Talents.Attack - 1
}
