package core

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/genshinsim/gcsim/internal/eventlog"
	"github.com/genshinsim/gcsim/internal/player"
	"github.com/genshinsim/gcsim/internal/tmpl/event"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

const (
	// MaxStam            = 240
	StamCDFrames       = 90
	JumpFrames         = 33
	DashFrames         = 24
	WalkFrames         = 1
	SwapCDFrames       = 60
	MaxTeamPlayerCount = 4
	DefaultTargetIndex = 1
)

type Core struct {
	coretype.Logger
	coretype.EventEmitter
	//control
	F     int            // current frame
	Flags coretype.Flags // global flags
	Rand  *rand.Rand
	Log   coretype.Logger

	//Player
	Player player.Player

	//track targets
	Targets     []Target
	TotalDamage float64 // keeps tracks of total damage dealt for the purpose of final results

	//last action taken by the sim
	LastAction coretype.ActionItem

	//tracks the current animation state
	state       coretype.AnimationState
	stateExpiry int

	//handlers
	Status     StatusHandler
	Energy     EnergyHandler
	Action     CommandHandler
	Queue      QueueHandler
	Combat     CombatHandler
	Tasks      TaskHandler
	Constructs ConstructHandler
	Shields    ShieldHandler
	Health     HealthHandler
}

func New() *Core {
	// var err error
	c := &Core{}

	c.Logger = eventlog.NewCtrl(&c.F, 0)
	c.EventEmitter = event.NewCtrl(c)

	c.Flags.Custom = make(map[string]int)
	//make a default nil writer
	c.Log = &eventlog.NilLogger{}

	return c
}

func (c *Core) Init() {
	c.Player.Init()
	c.Emit(coretype.OnInitialize)
}

func (c *Core) AddChar(v coretype.CharacterProfile) (coretype.Character, error) {
	f, ok := charMap[v.Base.Key]
	if !ok {
		return nil, fmt.Errorf("invalid character: %v", v.Base.Key.String())
	}
	char, err := f(c, v)
	if err != nil {
		return nil, err
	}
	c.Player.AddChar(char)

	wf, ok := weaponMap[v.Weapon.Name]
	if !ok {
		return nil, fmt.Errorf("unrecognized weapon %v for character %v", v.Weapon.Name, v.Base.Key.String())
	}
	wk := wf(char, c, v.Weapon.Refine, v.Weapon.Params)
	char.SetWeaponKey(wk)

	//add set bonus
	total := 0
	for key, count := range v.Sets {
		total += count
		f, ok := setMap[key]
		if ok {
			f(char, c, count, v.SetParams[key])
		} else {
			return nil, fmt.Errorf("character %v has unrecognized artifact: %v", v.Base.Name, key)
		}
	}
	if total > 5 {
		return nil, fmt.Errorf("total set count cannot exceed 5, got %v", total)
	}

	err = char.CalcBaseStats()

	return char, err
}

func (c *Core) CharByName(key coretype.CharKey) (coretype.Character, bool) {
	return c.Player.CharByName(key)
}

func (c *Core) Swap(next coretype.CharKey) int {
	return c.Player.Swap(next)
}

func (c *Core) AnimationCancelDelay(next coretype.ActionType, p map[string]int) int {
	//if last action is jump, dash, swap,
	switch c.LastAction.Typ {
	case coretype.ActionSwap:
		fallthrough
	case coretype.ActionDash:
		fallthrough
	case coretype.ActionJump:
		return 0
	}
	//other wise check with the current character
	return c.Player.Chars[c.Player.ActiveChar].ActionInterruptableDelay(next, p)
}

func (c *Core) UserCustomDelay() int {
	d := 0
	switch c.LastAction.Typ {
	case coretype.ActionSkill:
		d = c.Flags.Delays.Skill
	case coretype.ActionBurst:
		d = c.Flags.Delays.Burst
	case coretype.ActionAttack:
		d = c.Flags.Delays.Attack
	case coretype.ActionCharge:
		d = c.Flags.Delays.Charge
	case coretype.ActionDash:
		d = c.Flags.Delays.Dash
	case coretype.ActionJump:
		d = c.Flags.Delays.Jump
	case coretype.ActionSwap:
		d = c.Flags.Delays.Swap
	case coretype.ActionAim:
		d = c.Flags.Delays.Aim
	}
	return c.LastAction.Param["delay"] + d
}

func (c *Core) SetCustomFlag(key string, val int) {
	c.Flags.Custom[key] = val
}

func (c *Core) GetCustomFlag(key string) (int, bool) {
	val, ok := c.Flags.Custom[key]
	return val, ok
}

func (c *Core) Skip(n int) {
	for i := 0; i < n; i++ {
		c.Tick()
	}
}

func (c *Core) Tick() {
	//increment frame count
	c.F++
	//tick auras
	for _, t := range c.Targets {
		if t == nil {
			log.Print("unexpected nil target?")
			log.Println(c.Targets)
		}
		t.Tick()
	}
	//tick shields
	c.Shields.Tick()
	//tick constructs
	c.Constructs.Tick()

	c.Player.Tick()
	//run queued tasks
	c.Tasks.Run()

}
