package player

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core"
)

//ErrActionNotReady is returned if the requested action is not ready; this could be
//due to any of the following:
//	- Insufficient energy (burst only)
//	- Ability on cooldown
//	- Player currently in animation
var ErrActionNotReady = errors.New("action is not ready yet; cannot be executed")

//Exec mirrors the idea of the in game buttons where you can press the button but
//it may be greyed out. If grey'd out it will return ErrActionNotReady. Otherwise
//if action was executed successfully then it will return 0
//
//The function takes 2 params:
//	- ActionType
//	- Param
//
//Just like in game this will always try and execute on the currently active character
func (p *Player) Exec(t core.ActionType, param map[string]int) error {
	char := p.core.Chars[p.core.ActiveChar]

	if !char.ActionReady(t, param) {
		return ErrActionNotReady
	}

	switch t {
	case core.ActionSkill:
	case core.ActionBurst:
	case core.ActionAttack:
	case core.ActionCharge:
	case core.ActionHighPlunge:
	case core.ActionLowPlunge:
		case
	}

	return nil
}
