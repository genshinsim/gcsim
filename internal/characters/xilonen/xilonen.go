package xilonen

import (
	"errors"
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	c2buffsInit()
	core.RegisterCharFunc(keys.Xilonen, NewChar)
}

type char struct {
	*tmpl.Character

	nightsoulState    *nightsoul.State
	nightsoulSrc      int
	sampleSrc         int
	samplersConverted int
	shredElements     map[attributes.Element]bool
	c6activated       bool
	samplersActivated bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	c.shredElements = map[attributes.Element]bool{}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.nightsoulState = nightsoul.New(c.Core, c.CharWrapper)
	samplers := make([]attributes.Element, 4) // four samplers, one is herself but will be skipped
	for i := 0; i < 4; i++ {
		samplers[i] = attributes.Geo
	}

	c.samplersConverted = 0
	msg := c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, "converting samplers")
	for _, ch := range c.Core.Player.Chars() {
		if ch.Index == c.Index {
			// skip Xilonen herself
			continue
		}
		msg.Write(fmt.Sprintf("%v", ch.Base.Key.String()), ch.Base.Element.String())
		switch ch.Base.Element {
		case attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro:
			samplers[ch.Index] = ch.Base.Element
			c.samplersConverted++
		}
	}

	for i, ele := range samplers {
		if i == c.Index {
			// skip Xilonen herself
			continue
		}
		c.shredElements[ele] = true
	}

	c.a1()
	c.a4()

	c.c2()
	c.c4Init()

	c.c6Stam()

	c.onExitField()
	return nil
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}

		if !c.nightsoulState.HasBlessing() {
			return false
		}

		c.exitNightsoul()
		c.DeleteStatus(c6key)
		return false
	}, "xilonen-exit")
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	if c.nightsoulState.HasBlessing() {
		// cannot CA in nightsoul blessing
		if a == action.ActionCharge {
			return errors.New("xilonen cannot charge in nightsoul blessing")
		}
	}

	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 12
	default:
		return 4
	}
}
