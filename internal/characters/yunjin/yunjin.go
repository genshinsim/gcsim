package yunjin

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Yunjin, NewChar)
}

type char struct {
	*tmpl.Character
	partyElementalTypes int
	c4bonus             []float64
	c6bonus             []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.partyElementalTypes = 0

	w.Character = &c

	return nil
}

// Occurs after all characters are loaded, so getPartyElementalTypeCounts works properly
func (c *char) Init() error {
	c.a4Init()
	c.burstProc()
	c.c4()
	if c.Base.Cons >= 6 {
		c.c6bonus = make([]float64, attributes.EndStatType)
		c.c6bonus[attributes.AtkSpd] = .12
	}

	return nil
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 12
	}
	return c.AnimationStartDelay(k)
}
