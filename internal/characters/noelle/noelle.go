package noelle

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func init() {
	core.RegisterCharFunc(keys.Noelle, NewChar)
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
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4

	c.a2()

	return &c, nil
}

/**

a2: shielding if fall below hp threshold, not implemented

a4: every 4 hit decrease breastplate cd by 1; implement as hook

c1: 100% healing, not implemented

c2: decrease stam consumption, to be implemented

c4: explodes for 400% when expired or destroyed; how to implement expired?

c6: sweeping time increase additional 50%; add 1s up to 10s everytime opponent killed (NOT IMPLEMENTED, NOTHING DIES)

**/

func (c *char) a2() {
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
			Abil:       "A2 Shield",
			AttackTag:  core.AttackTagNone,
		}
		snap := c.Snapshot(&ai)

		//add shield
		x := snap.BaseDef*(1+snap.Stats[core.DEFP]) + snap.Stats[core.DEF]
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldNoelleA2,
			HP:         4 * x,
			Ele:        core.Cryo,
			Expires:    c.Core.F + 1200, //20 sec
		})
		return false
	}, "noelle-a2")
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Core.Status.Duration("noelleq") > 0 {

		x := c.Base.Def*(1+ds.Stats[core.DEFP]) + ds.Stats[core.DEF]
		mult := defconv[c.TalentLvlBurst()]
		if c.Base.Cons == 6 {
			mult += 0.5
		}
		fa := mult * x
		c.Core.Log.Debugw("noelle burst", "frame", c.Core.F, "event", core.LogSnapshotEvent, "total def", x, "atk added", fa, "mult", mult)

		ds.Stats[core.ATK] += fa
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
