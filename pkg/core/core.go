// Package core provides core functionality for a simulation:
//   - combat
//   - tasks
//   - event handling
//   - logging
//   - constructs (really should be just generic objects?)
//   - status
package core

import (
	"fmt"
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/status"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type Core struct {
	F     int
	Flags Flags
	Seed  int64
	Rand  *rand.Rand
	// various functionalities of core
	Log        glog.Logger    // we use an interface here so that we can pass in a nil logger for all except 1 run
	Events     *event.Handler // track events: subscribe/unsubscribe/emit
	Status     *status.Handler
	Tasks      *task.Handler
	Combat     *combat.Handler
	Constructs *construct.Handler
	Player     *player.Handler
}

type Flags struct {
	LogDebug     bool // Used to determine logging level
	DamageMode   bool // for hp mode
	DefHalt      bool // for hitlag
	EnableHitlag bool // hitlag enabled
	ErCalc       bool
	Custom       map[string]int
}

type Reactable interface {
	React(a *combat.AttackEvent)
	AuraContains(e ...attributes.Element) bool
	Tick()
}

// type Enemy interface {
// 	AddResistMod(key string, dur int, ele attributes.Element, val float64)
// 	DeleteResistMod(key string)
// 	ResistModIsActive(key string) bool
// 	AddDefMod(key string, dur int, val float64)
// 	DeleteDefMod(key string)
// 	DefModIsActive(key string) bool
// }

const MaxTeamSize = 4

type Opt struct {
	Seed         int64
	Debug        bool
	EnableHitlag bool
	DefHalt      bool
	DamageMode   bool
	ErCalc       bool
	Delays       info.Delays
}

func New(opt Opt) (*Core, error) {
	c := &Core{}
	c.Seed = opt.Seed
	c.Rand = rand.New(rand.NewSource(opt.Seed))
	c.Flags.Custom = make(map[string]int)
	if opt.Debug {
		c.Log = glog.New(&c.F, 500)
		c.Flags.LogDebug = true
	} else {
		c.Log = &glog.NilLogger{}
	}

	c.Flags.DamageMode = opt.DamageMode
	c.Flags.DefHalt = opt.DefHalt
	c.Flags.EnableHitlag = opt.EnableHitlag
	c.Flags.ErCalc = opt.ErCalc
	c.Events = event.New()
	c.Status = status.New(&c.F, c.Log)
	c.Tasks = task.New(&c.F)
	c.Constructs = construct.New(&c.F, c.Log)
	c.Player = player.New(
		player.Opt{
			F:            &c.F,
			Delays:       opt.Delays,
			Log:          c.Log,
			Events:       c.Events,
			Tasks:        c.Tasks,
			Debug:        opt.Debug,
			EnableHitlag: opt.EnableHitlag,
		},
	)
	c.Combat = combat.New(combat.Opt{
		Events:       c.Events,
		Team:         c.Player,
		Rand:         c.Rand,
		Debug:        c.Flags.LogDebug,
		Log:          c.Log,
		DamageMode:   c.Flags.DamageMode,
		DefHalt:      c.Flags.DefHalt,
		EnableHitlag: c.Flags.EnableHitlag,
		Tasks:        c.Tasks,
	})

	return c, nil
}

func (c *Core) Init() error {
	var err error
	// setup list
	//	- resonance
	//	- on hit energy
	//	- base stats
	//	- char inits
	//	- init call backs
	c.SetupOnNormalHitEnergy()
	err = c.Player.InitializeTeam()
	if err != nil {
		return err
	}
	c.Events.Emit(event.OnInitialize)
	return nil
}

func (c *Core) Tick() error {
	// things to tick:
	//	- targets
	//	- constructs
	//	- player (stamina, swap, animation, etc...)
	//		- character
	//		- shields
	//		- animation
	//		- stamina
	//		- swap
	//	- tasks
	//TODO: check for errors here?
	c.Combat.Tick()
	c.Constructs.Tick()
	c.Player.Tick()
	c.Tasks.Run()
	return nil
}

func (c *Core) AddChar(p info.CharacterProfile) (int, error) {
	var err error

	// initialize character
	char, err := character.New(p, &c.F, c.Flags.LogDebug, c.Log, c.Events, c.Tasks)
	if err != nil {
		return -1, err
	}

	f, ok := charMap[p.Base.Key]
	if !ok {
		return -1, fmt.Errorf("invalid character: %v", p.Base.Key.String())
	}
	err = f(c, char, p)
	if err != nil {
		return -1, err
	}
	index := c.Player.AddChar(char)

	// get starting hp
	char.StartHP = -1
	if hp, ok := p.Params["start_hp"]; ok {
		char.StartHP = hp
	}

	// set the energy
	char.Energy = char.EnergyMax
	if e, ok := p.Params["start_energy"]; ok {
		char.Energy = float64(e)
		// some sanity check in case user decide to set energy = 10000000
		if char.Energy > char.EnergyMax {
			char.Energy = char.EnergyMax
		}
	}

	// initialize weapon
	wf, ok := weaponMap[p.Weapon.Key]
	if !ok {
		return -1, fmt.Errorf("unrecognized weapon %v for character %v", p.Weapon.Key, p.Base.Key.String())
	}
	weap, err := wf(c, char, p.Weapon)
	if err != nil {
		return -1, err
	}
	char.SetWeapon(weap)

	// set bonus
	total := 0
	for key, count := range p.Sets {
		total += count
		af, ok := setMap[key]
		if ok {
			s, err := af(c, char, count, p.SetParams[key])
			if err != nil {
				return -1, err
			}
			char.SetArtifactSet(key, s)
		} else {
			return -1, fmt.Errorf("character %v has unrecognized artifact: %v", p.Base.Key.String(), key)
		}
	}
	//TODO: this should be handled by parser
	if total > 5 {
		return -1, fmt.Errorf("total set count cannot exceed 5, got %v", total)
	}

	return index, nil
}
