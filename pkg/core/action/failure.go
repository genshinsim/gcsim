package action

import (
	"encoding/json"
	"errors"
	"strings"
)

type ActionFailure int

const (
	NoFailure ActionFailure = iota
	SwapCD
	SkillCD
	BurstCD
	InsufficientEnergy
	InsufficientStamina
	CharacterDeceased // TODO: need chars to die first
	DashCD
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
}

func (e ActionFailure) String() string {
	return failureString[e]
}

func (e ActionFailure) MarshalJSON() ([]byte, error) {
	return json.Marshal(failureString[e])
}

func (e *ActionFailure) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range failureString {
		if v == s {
			*e = ActionFailure(i)
			return nil
		}
	}
	return errors.New("unrecognized ActionFailure")
}
