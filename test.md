# Table of Content

- [Sim structure](#sim-structure)
- [Core](#core)
  - [Shields and Constructs](#shields-and-constructs)
- [Monster](#monster)
  - [Hitbox](#hitbox)
  - [Damage](#damage)
  - [Auras](#auras)
  - [Reactions](#reactions)
- [Characters](#characters)
- [Weapons](#weapons)
- [Artifacts](#artifacts)
- [Parse](#parse)
- [Implementing New Character](#implementing-new-character)
- [Implementing New Weapon](#implementing-new-weapon)
- [Implementing New Artifacts](#implementing-new-artifacts)

# Sim structure

The sim is structured sort of a "hub and spoke" style. At the center is the hub, or the `core`. The `core` connects together various components of the sim including:

- targets: implements damage calculation and reactions
- characters: implements individual character logic
- shields: handles tracking shields and calculating shield damage taken
- constructs: handles tracking constructs
- queue: handles determining the next action to execute
- action: handles action execution; i.e. figure out which character function to call
- tasks: handles queuing and execution of future tasks
- health: handles player taking damage
- events: handles events and event callbacks
- energy: handles distributing energy particles
- combat: handles applying damage to the targets (not a strictly required component but is here just for separation)
- status: tracks various status, usually for buffs/effects uptime but can be used for anything with a duration

In this way, each component can talk to each other component by going through the `core`. For example, a character may wish to find out from the construct component how many active constructs there are (i.e. for Zhongli E ticks).

The most common is probably the use of statuses. For example, the following tracks Bennett's field uptime:

```go
c.Core.Status.AddStatus("btburst", 720)
```

The `core` specifies (via interfaces) the methods that must be implemented by each of the components. For example, the following is the interface that all targets must implement:

```go
type Target interface {
	Index() int
	SetIndex(ind int) //update the current index
	MaxHP() float64
	HP() float64
	//aura/reactions
	AuraType() EleType
	AuraContains(e ...EleType) bool
	Tick() //this should happen first before task ticks

	//attacks
	Attack(ds *Snapshot) (float64, bool)

	AddDefMod(key string, val float64, dur int)
	AddResMod(key string, val ResistMod)
	RemoveResMod(key string)
	RemoveDefMod(key string)
	HasDefMod(key string) bool
	HasResMod(key string) bool

	Delete() //gracefully deference everything so that it can be gc'd
}
```

The purpose for designing it this way is so that each individual component can be overwritten with a custom implementation. While realistically, there's no need for multiple implementation to run the sim (in fact there is a default implementation for each component), the reason why it's designed like this is so that for testing purposes, you may need to overwrite certain component implementation. For example, if you are testing a character, you may wish to overwrite the shield component in order to collect additional information that is not tracked directly by the sim.

# Core

The `core` [package](https://github.com/genshinsim/gcsim/tree/main/pkg/core) contains a `Core` structure which is the hub in the hub and spoke structure described above.

```go
type Core struct {
	//control
	F     int   // current frame
	Flags Flags // global flags
	Rand  *rand.Rand
	Log   *zap.SugaredLogger

	//core data
	Stam   float64
	SwapCD int

	//core stuff
	queue        []ActionItem
	stamModifier []func(a ActionType) (float64, bool)
	lastStamUse  int

	//track characters
	ActiveChar     int            // index of currently active char
	ActiveDuration int            // duration in frames that the current char has been on field for
	Chars          []Character    // array holding all the characters on the team
	charPos        map[string]int // map of character string name to their index (for quick lookup by name)

	//track targets
	Targets     []Target
	TotalDamage float64 // keeps tracks of total damage dealt for the purpose of final results

	//last action taken by the sim
	LastAction ActionItem

	//tracks the current animation state
	state       AnimationState
	stateExpiry int

	//handlers
	Status     StatusHandler
	Energy     EnergyHandler
	Action     ActionHandler
	Queue      QueueHandler
	Combat     CombatHandler
	Tasks      TaskHandler
	Constructs ConstructHandler
	Shields    ShieldHandler
	Health     HealthHandler
	Events     EventHandler
}
```

In addition to linking together all the component (or the handlers in the struct), it also keeps track of global variables that are shared among the various components. The most commonly used being `F`, which is the current frame.

The core also contains default implementation for the following components:

- [status](https://github.com/genshinsim/gcsim/blob/main/pkg/core/status.go)
- [energy](https://github.com/genshinsim/gcsim/blob/main/pkg/core/energy.go)
- [action](https://github.com/genshinsim/gcsim/blob/main/pkg/core/action.go)
- [queue](https://github.com/genshinsim/gcsim/blob/main/pkg/core/queue.go)
- [combat](https://github.com/genshinsim/gcsim/blob/main/pkg/core/combat.go)
- [tasks](https://github.com/genshinsim/gcsim/blob/main/pkg/core/tasks.go)
- [constructs](https://github.com/genshinsim/gcsim/blob/main/pkg/core/construct.go)
- [shields](https://github.com/genshinsim/gcsim/blob/main/pkg/core/shield.go)
- [health](https://github.com/genshinsim/gcsim/blob/main/pkg/core/health.go)
- [events](https://github.com/genshinsim/gcsim/blob/main/pkg/core/events.go)

`Target` and `Character` are a little bit special and are handled in their own packages (which will be covered in a later section). Ideally these default implementation should probably be split off into their own packages instead of being all in the `core` package...

## Shields and Constructs

Special note re. `shields` and `constructs`. `ShieldHandler` and `ConstructHandler` are the components that handles the tracking etc... of shields and construct. There is an additional interface definition for the shield and constructs themselves as follows:

```go
type Shield interface {
	Key() int
	Type() ShieldType
	OnDamage(dmg float64, ele EleType, bonus float64) (float64, bool) //return dmg taken and shield stays
	OnExpire()
	OnOverwrite()
	Expiry() int
	CurrentHP() float64
	Element() EleType
	Desc() string
}

type Construct interface {
	OnDestruct()
	Key() int
	Type() GeoConstructType
	Expiry() int
	IsLimited() bool
	Count() int
}
```

This is so that each character can implement their own logic for shields and constructs as some of them may have special effects (i.e. Noelle's shield doing damage on expiry, or Geo MC's rock doing damage on expiry)

# Monster

The `monster` package handles all the logic relating to:

- hitbox resolution (although not currently implemented)
- damage calculation
- ICD
- aura tracking and reactions
- target resistance and defense mods

All of this is implemented in a `Target` struct, which implements the `core.Target` interface. Thus multi target is simply having multiple copies of this `Target` struct. We'll call this the "target" for simplicity.

## Hitbox

When an attack is generated, whether or not a target will be hit/damaged is resolved by each target independently. Each of the attack [snapshot](https://github.com/genshinsim/gcsim/blob/main/pkg/core/snapshot.go) contains the information necessary for the target to determine if it will be hit or not.

For now the implementation is relatively simple. There is a `Targets` field in each snapshot. If this field is equal to the index of the current target or if this field is equal to -1 (representing all targets), then the current target will take damage.

In the future, this implementation can be changed to include 2D geometry.

## Damage

## Auras

## Reactions

# Characters

# Weapons

# Artifacts

# Parse

The parse [package](https://github.com/genshinsim/gcsim/tree/main/pkg/parse) contains the necessary code to lex/parse the custom config file syntax into the config data structure that's used by the core. The core logic is based on Rob Pike's [talk](https://talks.golang.org/2011/lex.slide#1) as well as go's template parsing [implementation](https://cs.opensource.google/go/go/+/refs/tags/go1.17:src/text/template/parse/)

# Implementing New Character

# Implementing New Weapon

# Implementing New Artifacts
