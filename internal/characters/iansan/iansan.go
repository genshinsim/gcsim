package iansan

import (
	"github.com/genshinsim/gcsim/internal/template/character/basicimport"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct{ *basicimport.Character }

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := &char{Character: basicimport.New(s, w, generatedProfile)}
	w.Character = c
	return nil
}
