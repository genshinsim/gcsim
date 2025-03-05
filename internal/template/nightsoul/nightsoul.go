package nightsoul

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const NightsoulBlessingStatus = "nightsoul-blessing"

type State struct {
	char            *character.CharWrapper
	c               *core.Core
	nightsoulPoints float64
	ExitStateF      int
	MaxPoints       float64
}

func New(c *core.Core, char *character.CharWrapper) *State {
	t := &State{
		char:      char,
		c:         c,
		MaxPoints: -1.0, // no limits
	}
	return t
}

func (s *State) Duration() int {
	return s.char.StatusDuration(NightsoulBlessingStatus)
}

// Change the current duration of the NS status if applicable, and apply cb on expiry.
// If NS is not currently active, do nothing.
func (s *State) SetNightsoulExitTimer(duration int, cb func()) {
	if !s.char.StatusIsActive(NightsoulBlessingStatus) {
		return
	}

	s.char.AddStatus(NightsoulBlessingStatus, duration, true)
	src := s.c.F + duration
	s.ExitStateF = src
	s.char.QueueCharTask(func() {
		if s.ExitStateF != src {
			return
		}
		cb()
	}, duration)
}

// Enters NS blessing with specified points.
// If duration is not infinite, expire NS upon duration and optionally trigger a CB.
// CB will not be called if the duration is changed before expiry.
func (s *State) EnterTimedBlessing(amount float64, duration int, cb func()) {
	s.nightsoulPoints = amount
	s.char.AddStatus(NightsoulBlessingStatus, duration, true)

	if duration > 0 && cb != nil {
		s.SetNightsoulExitTimer(duration, cb)
	}
	s.c.Log.NewEvent("enter nightsoul blessing", glog.LogCharacterEvent, s.char.Index).
		Write("points", s.nightsoulPoints).
		Write("duration", duration)
}

func (s *State) EnterBlessing(amount float64) {
	s.EnterTimedBlessing(amount, -1, nil)
}

func (s *State) ExitBlessing() {
	s.ExitStateF = 0
	s.char.DeleteStatus(NightsoulBlessingStatus)
	s.c.Log.NewEvent("exit nightsoul blessing", glog.LogCharacterEvent, s.char.Index)
}

func (s *State) HasBlessing() bool {
	return s.char.StatusIsActive(NightsoulBlessingStatus)
}

func (s *State) GeneratePoints(amount float64) {
	prevPoints := s.nightsoulPoints
	s.nightsoulPoints += amount
	s.clampPoints()
	s.c.Events.Emit(event.OnNightsoulGenerate, s.char.Index, amount)
	s.c.Log.NewEvent("generate nightsoul points", glog.LogCharacterEvent, s.char.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", s.nightsoulPoints)
}

func (s *State) ConsumePoints(amount float64) {
	prevPoints := s.nightsoulPoints
	s.nightsoulPoints -= amount
	s.clampPoints()
	s.c.Events.Emit(event.OnNightsoulConsume, s.char.Index, amount)
	s.c.Log.NewEvent("consume nightsoul points", glog.LogCharacterEvent, s.char.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", s.nightsoulPoints)
}

func (s *State) ClearPoints() {
	amt := s.nightsoulPoints
	s.nightsoulPoints = 0
	s.c.Events.Emit(event.OnNightsoulConsume, s.char.Index, amt)
	s.c.Log.NewEvent("clear nightsoul points", glog.LogCharacterEvent, s.char.Index).
		Write("previous points", amt)
}

func (s *State) clampPoints() {
	if s.MaxPoints > 0 && s.nightsoulPoints > s.MaxPoints {
		s.nightsoulPoints = s.MaxPoints
	} else if s.nightsoulPoints < 0 {
		s.nightsoulPoints = 0
	}
}

func (s *State) Points() float64 {
	return s.nightsoulPoints
}

func (s *State) Condition(fields []string) (any, error) {
	switch fields[1] {
	case "state":
		return s.HasBlessing(), nil
	case "points":
		return s.Points(), nil
	default:
		return nil, fmt.Errorf("invalid nightsoul condition: %v", fields[1])
	}
}
