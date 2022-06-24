package kazuha

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Kazuha, NewChar)
}

type char struct {
	*tmpl.Character
	a1Ele               attributes.Element
	qInfuse             attributes.Element
	c6Active            int
	infuseCheckLocation combat.AttackPattern
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Anemo
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassSword
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneInazuma

	c.infuseCheckLocation = combat.NewDefCircHit(1.5, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableObject)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	return nil
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.Base.Cons < 6 {
		return ds
	}

	if c.c6Active <= c.Core.F {
		return ds
	}

	//add 0.2% dmg for every EM
	ds.Stats[attributes.DmgP] += 0.002 * ds.Stats[attributes.EM]
	c.Core.Log.NewEvent("c6 adding dmg", glog.LogCharacterEvent, c.Index, "em", ds.Stats[attributes.EM], "final", ds.Stats[attributes.DmgP])
	return ds
}
