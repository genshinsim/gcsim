package zhongli

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Zhongli, NewChar)
}

type char struct {
	*tmpl.Character
	steleSnapshot combat.AttackEvent
	maxStele      int
	steleCount    int
}

// TODO: need to clean up zhongli code still
func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	c.maxStele = 1
	if c.Base.Cons >= 1 {
		c.maxStele = 2
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	return nil
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "shielded":
		return c.Tags["shielded"], nil
	default:
		return c.Character.Condition(fields)
	}
}
