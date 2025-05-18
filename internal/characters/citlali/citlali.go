package citlali

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Citlali, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState   *nightsoul.State
	itzpapaSrc       int
	skillShield      *shd
	numStellarBlades int
	numC6Stacks      float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = 3
	c.SkillCon = 3
	c.BurstCon = 5
	c.HasArkhe = false

	w.Character = &c

	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 100

	c.itzpapaSrc = -1
	c.numC6Stacks = 0

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()

	c.c1()
	c.c2()
	c.c6()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "stellar-blade":
		return c.numStellarBlades, nil
	case "opal-fire":
		return c.StatusIsActive(opalFireStateKey), nil
	case "c6-stacks":
		return c.numC6Stacks, nil
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 13
	case model.AnimationYelanN0StartDelay:
		return 3
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
