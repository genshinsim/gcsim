package kirara

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Kirara, NewChar)
}

type char struct {
	*tmpl.Character
	a1Stacks    int
	cardamoms   int
	mineSnap    combat.Snapshot
	minePattern combat.AttackPattern
	c6Buff      []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Ascension >= 4 {
		c.a4()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons >= 6 {
		c.c6Buff = make([]float64, attributes.EndStatType)
		c.c6Buff[attributes.PyroP] = 0.12
		c.c6Buff[attributes.HydroP] = 0.12
		c.c6Buff[attributes.CryoP] = 0.12
		c.c6Buff[attributes.ElectroP] = 0.12
		c.c6Buff[attributes.AnemoP] = 0.12
		c.c6Buff[attributes.GeoP] = 0.12
		c.c6Buff[attributes.PhyP] = 0.12
		c.c6Buff[attributes.DendroP] = 0.12
	}

	return nil
}
