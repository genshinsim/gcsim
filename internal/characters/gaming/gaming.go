package gaming

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Gaming, NewChar)
}

type char struct {
	*tmpl.Character
	specialPlungeRadius float64
	manChaiWalkBack     int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	c.manChaiWalkBack = 92 // default assumption

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.c2()
	c.c6()
	c.onExitField()
	return nil
}
