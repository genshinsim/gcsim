package sandrone

import (
	"errors"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const maxDecode = 1000

type fagioState int

const (
	stateIdle fagioState = iota
	stateDecoding
	stateOverdriveCharge
	stateOverdriveIdle
)

type char struct {
	*tmpl.Character

	// decode values are stored as 10x the actual value to avoid floating point errors
	decode               int
	currSweepingFireTick int
	caLength             int
	caToBeams            bool
	caMaxBeams           int
	caBeamsRemaining     int
	caSrc                int
	decodeReduceSrc      int
	currFagio            fagioState
	// decode values are stored as 10x the actual value to avoid floating point errors
	a1DecreasedDecode int
	a1RefinedTactics  int
	c2stacks          int
	c6stacks          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.NormalCon = 3
	c.NormalHitNum = normalHitNum

	w.Character = &c

	c.caMaxBeams = 3

	return nil
}

func (c *char) Init() error {
	c.a4Init()
	c.stellarInit()
	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	if c.Core.Player.CurrentState() == action.SkillState {
		if a == action.ActionJump {
			return errors.New("cannot jump during skill")
		}

		if a == action.ActionWalk {
			return errors.New("cannot walk during skill")
		}
	}

	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 24
	case info.AnimationYelanN0StartDelay:
		return 17
	}
	return c.Character.AnimationStartDelay(k)
}
