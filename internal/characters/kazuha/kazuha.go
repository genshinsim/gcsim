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
	infuseCheckLocation combat.AttackPattern
	c2buff              []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	c.infuseCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 1.5, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableObject)

	w.Character = &c

	return nil
}

const c6BuffKey = "kazuha-c6"

func (c *char) Init() error {
	c.a4()
	if c.Base.Cons >= 2 {
		c.c2buff = make([]float64, attributes.EndStatType)
		c.c2buff[attributes.EM] = 200
	}
	return nil
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.Base.Cons < 6 {
		return ds
	}
	if !c.StatusIsActive(c6BuffKey) {
		return ds
	}

	//add 0.2% dmg for every EM
	ds.Stats[attributes.DmgP] += 0.002 * ds.Stats[attributes.EM]
	c.Core.Log.NewEvent("c6 adding dmg", glog.LogCharacterEvent, c.Index).
		Write("em", ds.Stats[attributes.EM]).
		Write("final", ds.Stats[attributes.DmgP])
	return ds
}
