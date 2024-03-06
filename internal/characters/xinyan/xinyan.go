package xinyan

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type char struct {
	*tmpl.Character
	shieldLevel             int
	shieldLevel2Requirement int
	shieldLevel3Requirement int
	c2Buff                  []float64
	shieldTickSrc           int
}

func init() {
	core.RegisterCharFunc(keys.Xinyan, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	c.shieldLevel2Requirement = 2
	c.shieldLevel3Requirement = 3

	w.Character = &c

	c.shieldLevel = 1

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 28
	}
	return c.AnimationStartDelay(k)
}
