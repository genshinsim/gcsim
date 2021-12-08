package core

type EventType int

const (
	OnAttackWillLand EventType = iota //target, AttackEvent
	OnDamage                          //target, AttackEvent, amount, crit
	//reaction related
	// OnReactionOccured //target, AttackEvent
	// OnTransReaction   //target, AttackEvent
	// OnAmpReaction     //target, AttackEvent

	OnAuraDurabilityAdded    //target, ele, durability
	OnAuraDurabilityDepleted //target, ele
	// OnReaction               //target, AttackEvent, ReactionType
	ReactionEventStartDelim
	OnOverload       //target, AttackEvent
	OnSuperconduct   //target, AttackEvent
	OnMelt           //target, AttackEvent
	OnVaporize       //target, AttackEvent
	OnFrozen         //target, AttackEvent
	OnElectroCharged //target, AttackEvent
	OnSwirlHydro
	OnSwirlCryo
	OnSwirlElectro
	OnSwirlPyro
	OnCrystallizeHydro
	OnCrystallizeCryo
	OnCrystallizeElectro
	OnCrystallizePyro
	ReactionEventEndDelim
	//other stuff
	OnStamUse          //abil
	OnShielded         //shield
	OnCharacterSwap    //prev, next
	OnDash             //nil
	OnParticleReceived //particle
	OnTargetDied       //target
	OnCharacterHurt    //nil
	OnHeal             //src char, target character, amount
	//ability use
	PreSkill         //nil
	PostSkill        //nil, frames
	PreBurst         //nil
	PostBurst        //nil, frames
	PreAttack        //nil
	PostAttack       //nil, frames
	PreChargeAttack  //nil
	PostChargeAttack //nil, frames
	PrePlunge        //nil
	PostPlunge       //nil, frames
	PreAimShoot      //nil
	PostAimShoot     //nil, frames
	//sim stuff
	OnInitialize  //nil
	OnStateChange //prev, next
	EndEventTypes //elim
)

type EventHandler interface {
	Subscribe(e EventType, f EventHook, key string)
	Unsubscribe(e EventType, key string)
	Emit(e EventType, args ...interface{})
}

type EventHook func(args ...interface{}) bool

type EventCtrl struct {
	c      *Core
	events [][]ehook
}

type ehook struct {
	f   EventHook
	key string
	src int
}

func NewEventCtrl(c *Core) *EventCtrl {
	h := &EventCtrl{c: c}

	h.events = make([][]ehook, EndEventTypes)

	for i := range h.events {
		h.events[i] = make([]ehook, 0, 10)
	}

	return h
}

func (h *EventCtrl) Subscribe(e EventType, f EventHook, key string) {
	a := h.events[e]

	//check if override first
	ind := len(a)
	for i, v := range a {
		if v.key == key {
			ind = i
		}
	}
	if ind != 0 && ind != len(a) {
		h.c.Log.Debugw("hook added", "frame", h.c.F, "event", LogHookEvent, "overwrite", true, "key", key, "type", e)
		a[ind] = ehook{
			f:   f,
			key: key,
			src: h.c.F,
		}
	} else {
		a = append(a, ehook{
			f:   f,
			key: key,
			src: h.c.F,
		})
		h.c.Log.Debugw("hook added", "frame", h.c.F, "event", LogHookEvent, "overwrite", true, "key", key, "type", e)
	}
	h.events[e] = a
}

func (h *EventCtrl) Unsubscribe(e EventType, key string) {
	n := 0
	for _, v := range h.events[e] {
		if v.key != key {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}

func (h *EventCtrl) Emit(e EventType, args ...interface{}) {
	n := 0
	for i, v := range h.events[e] {
		if v.f(args...) {
			h.c.Log.Debugw("event hook ended", "frame", h.c.F, "event", LogHookEvent, "key", i, "src", v.src)
		} else {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}
