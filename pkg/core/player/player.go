// Package player contains player related tracking and functionalities:
// - tracking characters on the team
// - handling animations state
// - handling normal attack state
// - handling character stats and attributes
// - handling shielding
package player

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
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
	// handlers
	*animation.AnimationHandler
	Shields *shield.Handler
	infusion.Handler

	// tracking
	chars   []*character.CharWrapper
	active  int
	charPos map[keys.Char]int

	// stam
	Stam            float64
	LastStamUse     int
	stamPercentMods []stamPercentMod

	// airborne source
	airborne AirborneSource

	// swap
	SwapCD int

	// dash: dash fails iff lockout && on CD
	DashCDExpirationFrame int
	DashLockout           bool

	// last action
	LastAction struct {
		UsedAt int
		Type   action.Action
		Param  map[string]int
		Char   int
	}
}

type Opt struct {
	F            *int
	Log          glog.Logger
	Events       event.Eventter
	Tasks        task.Tasker
	Delays       info.Delays
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
	h.Handler = infusion.New(opt.F, opt.Log, opt.Debug)
	h.AnimationHandler = animation.New(opt.F, opt.Debug, opt.Log, opt.Events, opt.Tasks)
	return h
}

func (h *Handler) swap(to keys.Char) func() {
	return func() {
		prev := h.active
		h.active = h.charPos[to]

		// still have remaining frames left on dash CD, save in char for when they go on-field again
		if h.DashCDExpirationFrame > *h.F {
			h.chars[prev].RemainingDashCD = h.DashCDExpirationFrame - *h.F
			h.chars[prev].DashLockout = h.DashLockout
		}

		// set the new DashCDExpirationFrame and reset character remaining back to 0
		h.DashCDExpirationFrame = *h.F + h.chars[h.active].RemainingDashCD
		h.DashLockout = h.chars[h.active].DashLockout
		h.chars[h.active].RemainingDashCD = 0

		h.SwapCD = SwapCDFrames
		h.ResetAllNormalCounter()

		evt := h.Log.NewEvent("executed swap", glog.LogActionEvent, h.active).
			Write("action", "swap").
			Write("target", to.String())

		if h.chars[prev].RemainingDashCD > 0 {
			evt.Write("prev_dash_cd", h.chars[prev].RemainingDashCD).
				Write("prev_dash_lockout", h.chars[prev].DashLockout)
		}

		if h.DashCDExpirationFrame > *h.F {
			evt.Write("target_dash_cd", h.DashCDExpirationFrame-*h.F).
				Write("target_dash_expiry_frame", h.DashCDExpirationFrame).
				Write("target_dash_lockout", h.DashLockout)
		}

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
	// make sure we only apply hitlag if this character is on field
	if char != h.active {
		return
	}

	h.chars[char].ApplyHitlag(factor, dur)

	// also extend infusion
	//TODO: this is a really awkward place to apply this
	h.ExtendInfusion(char, factor, dur)

	// extend the dash cd by the hitlag extension amount
	if h.DashCDExpirationFrame > *h.F {
		ext := int(math.Ceil(dur * (1 - factor)))
		h.DashCDExpirationFrame += ext

		var evt glog.Event
		if h.DashLockout {
			evt = h.Log.NewEvent("dash cd hitlag extended", glog.LogHitlagEvent, char)
		} else {
			evt = h.Log.NewEvent("dash lockout evaluation hitlag extended", glog.LogHitlagEvent, char)
		}
		evt.Write("extension", ext).
			Write("expiry", h.DashCDExpirationFrame-*h.F).
			Write("expiry_frame", h.DashCDExpirationFrame).
			Write("lockout", h.DashLockout)
	}
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
	// loop again to initialize
	for i := range h.chars {
		err = h.chars[i].Init()
		if err != nil {
			return err
		}
		h.chars[i].Equip.Weapon.Init()
		for k := range h.chars[i].Equip.Sets {
			h.chars[i].Equip.Sets[k].Init()
		}
		// set each char's starting hp ratio
		switch {
		case h.chars[i].StartHP > 0 && h.chars[i].StartHPRatio > 0:
			h.chars[i].SetHPByRatio(float64(h.chars[i].StartHPRatio) / 100.0)
			h.chars[i].ModifyHPByAmount(float64(h.chars[i].StartHP))
		case h.chars[i].StartHP > 0:
			h.chars[i].SetHPByAmount(float64(h.chars[i].StartHP))
		case h.chars[i].StartHPRatio > 0:
			h.chars[i].SetHPByRatio(float64(h.chars[i].StartHPRatio) / 100.0)
		default:
			h.chars[i].SetHPByRatio(1)
		}
		h.Log.NewEvent("starting hp set", glog.LogCharacterEvent, i).
			Write("starting_hp_ratio", h.chars[i].CurrentHPRatio()).
			Write("starting_hp", h.chars[i].CurrentHP())
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
	// recover stamina
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
	h.AnimationHandler.Tick()
	for _, c := range h.chars {
		c.Tick()
	}
}

type AirborneSource int

const (
	Grounded AirborneSource = iota
	AirborneXiao
	AirborneVenti
	AirborneKazuha
	AirborneXianyun
	AirborneOroron
	TerminateAirborne
)

func (h *Handler) SetAirborne(src AirborneSource) error {
	if src < Grounded || src >= TerminateAirborne {
		// do nothing
		return fmt.Errorf("invalid airborne source: %v", src)
	}
	h.airborne = src
	return nil
}

func (h *Handler) Airborne() AirborneSource {
	return h.airborne
}

const (
	XianyunAirborneBuff = "xianyun-airborne-buff"
)
