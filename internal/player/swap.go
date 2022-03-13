package player

import "github.com/genshinsim/gcsim/pkg/coretype"

func (p *Player) Swap(next coretype.CharKey) int {
	prev := p.ActiveChar
	p.ActiveChar = p.CharPos[next]
	p.SwapCD = SwapCDFrames
	p.ResetAllNormalCounter()
	p.core.Emit(coretype.OnCharacterSwap, prev, p.ActiveChar)
	//this duration reset needs to be after the hook for spine to behave properly
	p.ActiveDuration = 0
	return 1
}
