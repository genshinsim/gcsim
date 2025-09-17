package skirk

import (
	"errors"
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Skirk, NewChar)
}

type char struct {
	*tmpl.Character
	serpentsSubtlety float64
	skillSrc         int
	prevNASkillState bool
	burstCount       int
	burstVoids       int
	voidRifts        RingQueue[int]
	c2Atk            []float64
	c6Stacks         RingQueue[int]
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 0
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	w.Character = &c

	ss, ok := p.Params["start_serpents_subtlety"]
	if !ok {
		ss = maxSerpentsSubtlety
	}
	ss = max(min(ss, maxSerpentsSubtlety), 0)
	c.serpentsSubtlety = float64(ss)

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.particleInit()
	c.BurstInit()
	c.a1Init()
	c.a4Init()
	c.talentPassiveInit()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) talentPassiveInit() {
	allCryoHydro := true
	hasCryo := false
	hasHydro := false

	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Hydro:
			hasHydro = true
		case attributes.Cryo:
			hasCryo = true
		default:
			allCryoHydro = false
		}
	}
	if !allCryoHydro {
		return
	}
	if !hasCryo {
		return
	}
	if !hasHydro {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.SetTag(keys.SkirkPassive, 1)
	}
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "serpents_subtlety":
		return c.serpentsSubtlety, nil
	case "void_rifts":
		filter := func(src int) bool {
			return src+a1Dur >= c.Core.F
		}
		return c.voidRifts.Count(filter), nil
	case "a4_stacks":
		return c.getA4Stacks(), nil
	case "c6_stacks":
		if c.Base.Cons < 6 {
			return 0, nil
		}
		filter := func(src int) bool {
			return src+c6Dur >= c.TimePassed
		}
		count := c.c6Stacks.Count(filter)
		return count, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	// TODO: Adjust this value based on if windup happened for the NA
	if k == info.AnimationXingqiuN0StartDelay {
		return 12
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	switch a {
	case action.ActionBurst:
		if !c.Core.Flags.IgnoreBurstEnergy && !c.StatusIsActive(skillKey) && c.serpentsSubtlety < 50 {
			return false, action.InsufficientEnergy
		}
	case action.ActionSkill:
		if c.StatusIsActive(skillKey) {
			return false, action.SkillCD
		}
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	// can use charge without attack beforehand unlike most of the other sword users
	if a == action.ActionCharge {
		return nil
	}
	if c.StatusIsActive(skillKey) && a == action.ActionSkill {
		return errors.New("skirk: cannot use skill in seven-phase flash")
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) AddSerpentsSubtlety(src string, e float64) {
	pre := c.serpentsSubtlety
	c.serpentsSubtlety += e
	c.serpentsSubtlety = min(max(c.serpentsSubtlety, 0), maxSerpentsSubtlety)

	c.Core.Log.NewEvent(fmt.Sprintf("+%.1f serpent's subtlety, next: %.1f", e, c.serpentsSubtlety), glog.LogEnergyEvent, c.Index()).
		Write("added", e).
		Write("pre_recovery", pre).
		Write("post_recovery", c.serpentsSubtlety).
		Write("source", src)
}

func (c *char) ReduceSerpentsSubtlety(src string, e float64) {
	pre := c.serpentsSubtlety
	c.serpentsSubtlety -= e
	c.serpentsSubtlety = min(max(c.serpentsSubtlety, 0), maxSerpentsSubtlety)

	c.Core.Log.NewEvent(fmt.Sprintf("-%.1f serpent's subtlety, next: %.1f", e, c.serpentsSubtlety), glog.LogEnergyEvent, c.Index()).
		Write("reduced", e).
		Write("pre", pre).
		Write("post", c.serpentsSubtlety).
		Write("source", src)
}

// Consumes SS after a specified delay. Not hitlag affected
func (c *char) ConsumeSerpentsSubtlety(delay int, src string) {
	if delay == 0 {
		c.Core.Log.NewEvent("draining serpent's subtlety", glog.LogEnergyEvent, c.Index()).
			Write("pre_drain", c.serpentsSubtlety).
			Write("post_drain", 0).
			Write("source", src)
		c.serpentsSubtlety = 0
		return
	}
	c.Core.Tasks.Add(func() {
		c.Core.Log.NewEvent("draining serpent's subtlety", glog.LogEnergyEvent, c.Index()).
			Write("pre_drain", c.serpentsSubtlety).
			Write("post_drain", 0).
			Write("source", src)
		c.serpentsSubtlety = 0
	}, delay)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...any) bool {
		if c.StatusIsActive(skillKey) {
			c.exitSkillState(c.skillSrc)
		}

		return false
	}, "skirk-exit")
}
