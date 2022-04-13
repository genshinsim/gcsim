package noelle

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Noelle, NewChar)
}

type char struct {
	*character.Tmpl
	shieldTimer int
	a4Counter   int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Geo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4

	c.InitCancelFrames()

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a1()
}

/**

a1: shielding if fall below hp threshold, not implemented

a4: every 4 hit decrease breastplate cd by 1; implement as hook

c2: decrease stam consumption, to be implemented

c4: explodes for 400% when expired or destroyed; how to implement expired?

c6: sweeping time increase additional 50%; add 1s up to 10s everytime opponent killed (NOT IMPLEMENTED, NOTHING DIES)

**/

func (c *char) a1() {
	icd := 0
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.F < icd {
			return false
		}
		char := c.Core.Chars[c.Core.ActiveChar]
		if char.HP()/char.MaxHP() >= 0.3 {
			return false
		}
		icd = c.Core.F + 3600
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "A1 Shield",
			AttackTag:  core.AttackTagNone,
		}
		snap := c.Snapshot(&ai)

		//add shield
		x := snap.BaseDef*(1+snap.Stats[core.DEFP]) + snap.Stats[core.DEF]
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldNoelleA1,
			Name:       "Noelle A1",
			HP:         4 * x,
			Ele:        core.Cryo,
			Expires:    c.Core.F + 1200, //20 sec
		})
		return false
	}, "noelle-a1")
}

// Noelle Geo infusion can't be overridden, so it must be a snapshot modification rather than a weapon infuse
func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Core.Status.Duration("noelleq") > 0 {
		//infusion to attacks only
		switch ai.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagPlunge:
		case core.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = core.Geo
	}
	return ds
}
