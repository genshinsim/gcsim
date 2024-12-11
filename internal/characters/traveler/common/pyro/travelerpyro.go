package pyro

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Traveler struct {
	*tmpl.Character
	nightsoulState        *nightsoul.State
	nightsoulSrc          int
	scorchingThresholdICD int
	gender                int
	c2ActivationsPerSkill int
}

func NewTraveler(s *core.Core, w *character.CharWrapper, p info.CharacterProfile, gender int) (*Traveler, error) {
	c := Traveler{
		gender: gender,
	}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 70
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 5

	common.TravelerBaseAtkIncrease(w, p)

	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 80
	c.c2ActivationsPerSkill = 0

	return &c, nil
}

func (c *Traveler) Init() error {
	c.scorchingThresholdOnDamage()
	c.a4()
	return nil
}

func (c *Traveler) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		if c.gender == 0 {
			return 8
		} else {
			return 7
		}
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *Traveler) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	default:
		return c.Character.Condition(fields)
	}
}
