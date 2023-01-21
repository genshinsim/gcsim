// Package player contains player related tracking and functionalities:
// - tracking characters on the team
// - handling animations state
// - handling normal attack state
// - handling character stats and attributes
// - handling shielding
package player

import (
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/animation"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/infusion"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

const (
	MaxStam      = 240
	StamCDFrames = 90
	SwapCDFrames = 60
)

type Handler struct {
	Opt
	//handlers
	*animation.AnimationHandler
	Shields *shield.Handler
	infusion.InfusionHandler

	//tracking
	chars   []*character.CharWrapper
	active  int
	charPos map[keys.Char]int

	//stam
	Stam            float64
	LastStamUse     int
	stamPercentMods []stamPercentMod

	//swap
	SwapCD int

	//last action
	LastAction struct {
		UsedAt int
		Type   action.Action
		Param  map[string]int
		Char   int
	}
}

type Delays struct {
	Skill  int `json:"skill"`
	Burst  int `json:"burst"`
	Attack int `json:"attack"`
	Charge int `json:"charge"`
	Aim    int `json:"aim"`
	Dash   int `json:"dash"`
	Jump   int `json:"jump"`
	Swap   int `json:"swap"`
}

type Opt struct {
	F            *int
	Log          glog.Logger
	Events       event.Eventter
	Tasks        task.Tasker
	Delays       Delays
	Debug        bool
	EnableHitlag bool
}

func New(opt Opt) *Handler {
	h := &Handler{
		chars:           make([]*character.CharWrapper, 0, 4),
		charPos:         make(map[keys.Char]int),
		stamPercentMods: make([]stamPercentMod, 0, 5),
		Opt:             opt,
		Stam:            MaxStam,
	}
	h.Shields = shield.New(opt.F, opt.Log, opt.Events)
	h.InfusionHandler = infusion.New(opt.F, opt.Log, opt.Debug)
	h.AnimationHandler = animation.New(opt.F, opt.Debug, opt.Log, opt.Events, opt.Tasks)
	return h
}

func (h *Handler) swap(to keys.Char) func() {
	return func() {
		prev := h.active
		h.active = h.charPos[to]
		h.Log.NewEvent("executed swap", glog.LogActionEvent, h.active).
			Write("action", "swap").
			Write("target", to.String())
		h.SwapCD = SwapCDFrames
		h.ResetAllNormalCounter()
		h.Events.Emit(event.OnCharacterSwap, prev, h.active)
	}
}

func (h *Handler) AddChar(char *character.CharWrapper) int {
	h.chars = append(h.chars, char)
	index := len(h.chars) - 1
	char.SetIndex(index)
	h.charPos[char.Base.Key] = index

	return index
}

func (h *Handler) ByIndex(i int) *character.CharWrapper {
	return h.chars[i]
}

func (h *Handler) CombatByIndex(i int) combat.Character {
	return h.chars[i]
}

func (h *Handler) ByKey(k keys.Char) (*character.CharWrapper, bool) {
	i, ok := h.charPos[k]
	if !ok {
		return nil, false
	}
	return h.chars[i], true
}

func (h *Handler) Chars() []*character.CharWrapper {
	return h.chars
}

func (h *Handler) Active() int {
	return h.active
}

func (h *Handler) ActiveChar() *character.CharWrapper {
	return h.chars[h.active]
}

// returns the char with the lowest HP
func (h *Handler) LowestHPChar() *character.CharWrapper {
	result := make([]*character.CharWrapper, 0, len(h.chars))

	// filter out dead characters
	for _, c := range h.chars {
		if c.HPCurrent <= 0 {
			continue
		}
		result = append(result, c)
	}

	// sort by HP
	sort.Slice(result, func(i, j int) bool {
		return result[i].HPCurrent < result[j].HPCurrent
	})

	return result[0]
}

func (h *Handler) CharIsActive(k keys.Char) bool {
	return h.charPos[k] == h.active
}

func (h *Handler) SetActive(i int) {
	h.active = i
}

func (h *Handler) Adjust(src string, char int, amt float64) {
	h.chars[char].AddEnergy(src, amt)
}

func (h *Handler) ResetAllNormalCounter() {
	for _, char := range h.chars {
		char.ResetNormalCounter()
	}
}

func (h *Handler) DistributeParticle(p character.Particle) {
	for i, char := range h.chars {
		char.ReceiveParticle(p, h.active == i, len(h.chars))
	}
	h.Events.Emit(event.OnParticleReceived, p)
}

func (h *Handler) AbilStamCost(i int, a action.Action, p map[string]int) float64 {
	// stam percent mods are negative
	// cap it to 100% stam decrease
	r := 1 + h.StamPercentMod(a)
	if r < 0 {
		r = 0
	}
	return r * h.chars[i].ActionStam(a, p)
}
func (h *Handler) RestoreStam(v float64) {
	h.Stam += v
	if h.Stam > MaxStam {
		h.Stam = MaxStam
	}
}

func (h *Handler) ApplyHitlag(char int, factor, dur float64) {
	//make sure we only apply hitlag if this character is on field
	if char != h.active {
		return
	}
	h.chars[char].ApplyHitlag(factor, dur)
	//also extend infusion
	//TODO: this is a really awkward place to apply this
	h.ExtendInfusion(char, factor, dur)
}

// InitializeTeam will set up resonance event hooks and calculate
// all character base stats
func (h *Handler) InitializeTeam() error {
	var err error
	for _, c := range h.chars {
		err = c.UpdateBaseStats()
		if err != nil {
			return err
		}
	}
	//loop again to initialize
	for i := range h.chars {
		err = h.chars[i].Init()
		if err != nil {
			return err
		}
		h.chars[i].Equip.Weapon.Init()
		for k := range h.chars[i].Equip.Sets {
			h.chars[i].Equip.Sets[k].Init()
		}
		//set each char's starting hp
		if h.chars[i].HPCurrent == -1 {
			h.chars[i].HPCurrent = h.chars[i].MaxHP()
		}
		h.Log.NewEvent("starting hp set", glog.LogCharacterEvent, i).
			Write("hp", h.chars[i].HPCurrent)
	}
	return nil
}

func (h *Handler) Tick() {
	//	- player (stamina, swap, animation, etc...)
	//		- character
	//		- shields
	//		- animation
	//		- stamina
	//		- swap
	//recover stamina
	if h.Stam < MaxStam && *h.F-h.LastStamUse > StamCDFrames {
		h.Stam += 25.0 / 60
		if h.Stam > MaxStam {
			h.Stam = MaxStam
		}
	}
	if h.SwapCD > 0 {
		h.SwapCD--
	}
	h.Shields.Tick()
	for _, c := range h.chars {
		c.Tick()
	}
	h.AnimationHandler.Tick()

}
