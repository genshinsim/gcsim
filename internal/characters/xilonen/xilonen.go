package xilonen

import (
	"errors"

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
	core.RegisterCharFunc(keys.Xilonen, NewChar)
}

type char struct {
	*tmpl.Character

	nightsoulState    *nightsoul.State
	nightsoulSrc      int
	sampleSrc         int
	c2Src             int
	exitStateSrc      int
	samplersConverted int
	shredElements     map[attributes.Element]bool
	skillLastStamF    int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.nightsoulState = nightsoul.New(c.Core, c.CharWrapper)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	c.shredElements = map[attributes.Element]bool{}

	w.Character = &c
	c.nightsoulState.MaxPoints = 90

	return nil
}

func (c *char) Init() error {
	for _, other := range c.Core.Player.Chars() {
		if other.Index == c.Index {
			// skip Xilonen herself
			continue
		}
		switch ele := other.Base.Element; ele {
		case attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro:
			c.samplersConverted++
			c.shredElements[ele] = true
		default:
			c.shredElements[attributes.Geo] = true
		}
	}
	if len(c.Core.Player.Chars()) < 4 {
		c.shredElements[attributes.Geo] = true
	}

	c.a1()
	c.a4()

	c.c2()
	c.c4Init()
	c.c6()

	c.onExitField()
	return nil
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev == c.Index {
			c.exitNightsoul()
		}
		return false
	}, "xilonen-exit")
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	if a == action.ActionCharge && c.canUseNightsoul() {
		return errors.New("xilonen cannot charge in nightsoul blessing")
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if a1 window is active is on-field
	if a == action.ActionSkill && c.StatusIsActive(skilRecastCD) {
		return false, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 12
	case model.AnimationYelanN0StartDelay:
		return 4
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if c.canUseNightsoul() {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	default:
		return c.Character.Condition(fields)
	}
}
