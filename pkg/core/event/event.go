package event

type Event int

const (
	OnAttackWillLand Event = iota //target, AttackEvent
	OnDamage                      //target, AttackEvent, amount, crit
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
	OnAggravate //target, AttackEvent
	OnSpread    //target, AttackEvent
	OnQuicken
	OnBloom //target, AttackEvent
	OnHyperbloom
	OnBurgeon
	OnBurning //target, AttackEvent
	ReactionEventEndDelim
	OnDendroCore //Gadget
	//other stuff
	OnStamUse          //abil
	OnShielded         //shield
	OnCharacterSwap    //prev, next
	OnParticleReceived //particle
	OnEnergyChange     //character_received_index, pre_energy, energy_change, src (post-energy available in character_received)
	OnTargetDied       //target, AttackEvent
	OnCharacterHit     //nil <- this is for when the character is going to get hit but might be shielded from dmg
	OnCharacterHurt    //amount
	OnHeal             //src char, target character, amount
	//ability use
	OnActionFailed //ActiveCharIndex, action.Action, param, action.ActionFailure
	OnActionExec   //ActiveCharIndex, action.Action, param
	OnSkill        //nil
	OnBurst        //nil
	OnAttack       //nil
	OnChargeAttack //nil
	OnPlunge       //nil
	OnAimShoot     //nil
	OnDash
	//sim stuff
	OnInitialize  //nil
	OnStateChange //prev, next
	OnEnemyAdded  //t
	OnTick
	EndEventTypes //elim
)

type Handler struct {
	events [][]ehook
}

type EventHook func(args ...interface{}) bool

type Eventter interface {
	Subscribe(e Event, f EventHook, key string)
	Unsubscribe(e Event, key string)
	Emit(e Event, args ...interface{})
}

type ehook struct {
	f   EventHook
	key string
}

func New() *Handler {
	h := &Handler{
		events: make([][]ehook, EndEventTypes),
	}

	for i := range h.events {
		h.events[i] = make([]ehook, 0, 10)
	}

	return h
}

func (h *Handler) Subscribe(e Event, f EventHook, key string) {
	a := h.events[e]

	evt := ehook{
		f:   f,
		key: key,
	}

	//check if override first
	ind := -1
	for i, v := range a {
		if v.key == key {
			ind = i
		}
	}
	if ind > -1 {
		a[ind] = evt
	} else {
		a = append(a, evt)
	}
	h.events[e] = a
}

func (h *Handler) Unsubscribe(e Event, key string) {
	n := 0
	for _, v := range h.events[e] {
		if v.key != key {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}

func (h *Handler) Emit(e Event, args ...interface{}) {
	n := 0
	for _, v := range h.events[e] {
		if !v.f(args...) {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}
