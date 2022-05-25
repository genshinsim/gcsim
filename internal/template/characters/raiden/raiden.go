package raiden

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
	burstCastF     int
	eyeICD         int
	stacksConsumed float64
	stacks         float64
	restoreICD     int
	restoreCount   int
	applyC4        bool
	c6Count        int
	c6ICD          int
}

func init() {
	core.RegisterCharFunc(core.Raiden, NewChar)
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
		e = 90
	}
	c.Energy = float64(e)
	c.EnergyMax = 90
	c.Weapon.Class = core.WeaponClassSpear
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5
	c.CharZone = core.ZoneInazuma

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()
	c.InitCancelFrames()

	c.eyeOnDamage()
	c.onBurstStackCount()
	c.onSwapClearBurst()

	if c.Base.Cons == 6 {
		c.c6()
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Core.Status.Duration("raidenburst") == 0 {
			return 25
		}
		return 20
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) Snapshot(a *core.AttackInfo) core.Snapshot {
	s := c.Tmpl.Snapshot(a)

	//a1 add dmg based on ER%
	excess := int(s.Stats[core.ER] / 0.01)

	s.Stats[core.ElectroP] += float64(excess) * 0.004 /// 0.4% extra dmg
	c.Core.Log.NewEvent("a4 adding electro dmg", core.LogCharacterEvent, c.Index, "stacks", excess, "final", s.Stats[core.ElectroP])
	//
	////infusion to normal/plunge/charge
	//switch ds.AttackTag {
	//case core.AttackTagNormal:
	//case core.AttackTagExtra:
	//case core.AttackTagPlunge:
	//default:
	//	return ds
	//}
	//if c.Core.Status.Duration("raidenburst") > 0 {
	//	ds.Element = core.Electro
	//}

	return s
}
