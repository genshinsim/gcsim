package player

import (
	"github.com/genshinsim/gcsim/pkg/coretype"
)

const (
	MaxStam            = 240
	StamCDFrames       = 90
	SwapCDFrames       = 60
	MaxTeamPlayerCount = 4
	DefaultTargetIndex = 1
)

type Core interface {
	coretype.Framer
	coretype.EventEmitter
	coretype.Logger
}

type Player struct {
	core Core
	//Characters
	ActiveChar     int                      // index of currently active char
	ActiveDuration int                      // duration in frames that the current char has been on field for
	Chars          []coretype.Character     // array holding all the characters on the team
	CharPos        map[coretype.CharKey]int // map of character string name to their index (for quick lookup by name)

	//Stamina
	stamModifier []stamMod
	lastStamUse  int
	Stam         float64

	//swap related
	SwapCD int

	//player also need to handle:
	//	energy
	//	shields
	//	health

	//healing
	healingBonus    []func(healedCharIndex int) float64 // Array that holds functions calculating incoming healing bonus
	damageReduction []func() (float64, bool)

	//shields
	shields         []coretype.Shield
	shieldBonusFunc []func() float64
}

func New(c Core) *Player {
	p := &Player{
		core:            c,
		CharPos:         make(map[coretype.CharKey]int),
		Stam:            MaxStam,
		stamModifier:    make([]stamMod, 0, 10),
		shields:         make([]coretype.Shield, 0, coretype.EndShieldType),
		shieldBonusFunc: make([]func() float64, 0, 10),
	}

	return p
}

func (p *Player) Init() {
	for _, char := range p.Chars {
		char.Init()
	}
}

func (p *Player) Tick() {
	for _, char := range p.Chars {
		char.Tick()
	}
	//recover stamina
	if p.Stam < MaxStam && p.core.F()-p.lastStamUse > StamCDFrames {
		p.Stam += 25.0 / 60
		if p.Stam > MaxStam {
			p.Stam = MaxStam
		}
	}
	//recover swap cd
	if p.SwapCD > 0 {
		p.SwapCD--
	}
	//update activeduration
	p.ActiveDuration++
}
