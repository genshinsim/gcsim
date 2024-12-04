package mavuika

// no ICD on E, 5 particles

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Mavuika, NewChar)
}

type char struct {
	*tmpl.Character
	// normals
	normalBCounter int

	// skill
	nightsoulState         *nightsoul.State
	nightsoulSrc           int
	flamestriderModeActive bool

	// burst
	fightingSpirit         float64
	maxFightingSpirit      float64
	consumedFightingSpirit float64

	// passives
	baseA4Buff float64
	a4Buff     []float64

	// cons
	fightingSpiritMult float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 0
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5
	c.HasArkhe = false

	w.Character = &c

	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 80

	c.flamestriderModeActive = false
	c.fightingSpirit = 200
	c.fightingSpiritMult = 1
	c.a4Buff = make([]float64, attributes.EndStatType)

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.burstInit()

	c.a1()

	c.c1()
	c.c2BaseIncrease()
	return nil
}

func (c *char) AdvanceNormalIndex() {
	if c.nightsoulState.HasBlessing() {
		c.normalBCounter++
		if c.normalBCounter == bikeHitNum {
			c.normalBCounter = 0
		}
		return
	}
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *char) ResetNormalCounter() {
	c.normalBCounter = 0
	c.Character.ResetNormalCounter()
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// if in the Nightsoul Blessing, can press E
	if a == action.ActionSkill && c.nightsoulState.HasBlessing() {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	case "fighting_spirit":
		return c.fightingSpirit, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if c.nightsoulState.HasBlessing() && c.flamestriderModeActive {
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

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.StatusIsActive(crucibleOfDeathAndLifeStatus) {
			c.DeleteStatus(crucibleOfDeathAndLifeStatus)
		}
		return false
	}, "mavuika-exit")
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if c.nightsoulState.HasBlessing() && c.flamestriderModeActive {
		return 0
	}
	return 50
}
