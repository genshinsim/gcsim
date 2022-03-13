package player

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/coretype"
)

//AddChar adds a new character to the player, returning the index for this
//character
func (p *Player) AddChar(c coretype.Character) (int, error) {
	if len(p.Chars) == MaxTeamPlayerCount {
		return -1, errors.New("number of characters cannot exceed 4")
	}
	p.Chars = append(p.Chars, c)
	i := len(p.Chars) - 1
	p.CharPos[c.Key()] = i
	c.SetIndex(i)

	return i, nil
}

func (p *Player) CharByName(key coretype.CharKey) (coretype.Character, bool) {
	pos, ok := p.CharPos[key]
	if !ok {
		return nil, false
	}
	return p.Chars[pos], true
}

func (p *Player) ResetAllNormalCounter() {
	for _, char := range p.Chars {
		char.ResetNormalCounter()
	}
}
