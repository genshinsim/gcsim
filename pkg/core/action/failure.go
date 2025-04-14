package action

import (
	"encoding/json"
	"errors"
	"strings"
)

type Failure int

const (
	NoFailure Failure = iota
	SwapCD
	SkillCD
	BurstCD
	InsufficientEnergy
	InsufficientStamina
	CharacterDeceased // TODO: need chars to die first
	DashCD
	TimeManip
)

var failureString = [...]string{
	"no_failure",
	"swap_cd",
	"skill_cd",
	"burst_cd",
	"insufficient_energy",
	"insufficient_stamina",
	"character_deceased",
	"dash_cd",
	"time_manip",
}

func (e Failure) String() string {
	return failureString[e]
}

func (e Failure) MarshalJSON() ([]byte, error) {
	return json.Marshal(failureString[e])
}

func (e *Failure) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range failureString {
		if v == s {
			*e = Failure(i)
			return nil
		}
	}
	return errors.New("unrecognized ActionFailure")
}
