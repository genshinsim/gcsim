package cryo

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	stellarConductText = " (Stellar-Conduct)"
	stellarSwirlText   = " (Stellar Swirl)"
)

type Traveler struct {
	*tmpl.Character
	gender          int
	skillSrc        int
	flostglowStacks int
	skillTravel     int
	trueMoonBuff    bool
	trueMoonStacks  int
}

func NewTraveler(s *core.Core, w *character.CharWrapper, p info.CharacterProfile, gender int) (*Traveler, error) {
	c := Traveler{
		gender: gender,
	}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Cryo
	c.EnergyMax = 60
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	common.TravelerStoryBuffs(w, p)
	trueMoonBuff, okTrueMoonBuff := p.Params[common.TrueMoonParam]

	if !okTrueMoonBuff {
		trueMoonBuff = 1
	}
	c.trueMoonBuff = trueMoonBuff > 0

	return &c, nil
}

func (c *Traveler) Init() error {
	c.stellarInit()
	c.trueMoonInit()
	c.a4Init()
	c.c1Init()
	c.c2Init()

	return nil
}

func (c *Traveler) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		if c.gender == 0 {
			return 8
		}
		return 7
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *Traveler) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "frostglow":
		return c.flostglowStacks, nil
	case "icepoint":
		return c.trueMoonStacks, nil
	default:
		return c.Character.Condition(fields)
	}
}
