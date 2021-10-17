package sara

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("sara", NewChar)
}

type char struct {
	*character.Tmpl
	a4LastProc int
	c1LastProc int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

// Handles Sara c6
func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.Element != core.Electro {
			return false
		}
		if c.Core.Status.Duration("sarabuff"+ds.Actor) <= 0 {
			return false
		}
		ds.Stats[core.CD] += .6
		return false
	}, fmt.Sprintf("sara-c6"))
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}
