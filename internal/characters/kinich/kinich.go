package kinich

import (
	"math"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Kinich, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState           *nightsoul.State
	nightsoulSrc             int
	ajawSrc                  int
	normalSCounter           int
	characterAngularPosition float64 // [0, 360)
	blindSpotAngularPosition float64 // [0, 360)
	particlesGenerated       bool
	durationExtended         bool
	c2AoeIncreased           bool
	scaleskiperAttackInfo    combat.AttackInfo
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	c.HasArkhe = false

	w.Character = &c

	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 20

	c.characterAngularPosition = 0.
	c.blindSpotAngularPosition = -1.

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()

	c.c1()
	c.c4()

	c.onExitField()
	return nil
}

func (c *char) AdvanceNormalIndex() {
	if c.nightsoulState.HasBlessing() {
		c.normalSCounter++
		if c.normalSCounter == skillHitNum {
			c.normalSCounter = 0
		}
		return
	}
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *char) ResetNormalCounter() {
	c.normalSCounter = 0
	c.Character.ResetNormalCounter()
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	case "blind_spot":
		if c.blindSpotAngularPosition == -1. {
			return 0, nil
		} else {
			diff := c.blindSpotAngularPosition - c.characterAngularPosition
			if diff > 180 {
				diff -= 360
			} else if diff < -180 {
				diff += 360
			}
			return diff / math.Abs(diff), nil
		}
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if c.nightsoulState.HasBlessing() {
		switch k {
		case model.AnimationXingqiuN0StartDelay:
			return 11
		default:
			return 9
		}
	}
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 11
	default:
		return 11
	}
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	if c.nightsoulState.HasBlessing() {
		// cannot CA, LP, HP in nightsoul blessing
		if a == action.ActionCharge {
			return player.ErrInvalidChargeAction
		}
		if a == action.ActionLowPlunge {
			return player.ErrInvalidChargeAction
		}
		if a == action.ActionHighPlunge {
			return player.ErrInvalidChargeAction
		}
	}

	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.nightsoulState.HasBlessing() {
			c.cancelNightsoul()
		}
		return false
	}, "kinich-exit")
}
