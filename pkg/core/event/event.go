package event

type Event int

const (
	OnEnemyHit     Event = iota //target, AttackEvent
	OnPlayerHit                 //target, AttackEvent
	OnGadgetHit                 //target, AttackEvent
	OnEnemyDamage               //target, AttackEvent, amount, crit
	OnGadgetDamage              //target, AttackEvent
	//reaction related
	// OnReactionOccured //target, AttackEvent
	// OnTransReaction   //target, AttackEvent
	// OnAmpReaction     //target, AttackEvent

	OnAuraDurabilityAdded    //target, ele, durability
	OnAuraDurabilityDepleted //target, ele
	// OnReaction               //target, AttackEvent, ReactionType
	ReactionEventStartDelim
	OnOverload           //target, AttackEvent
	OnSuperconduct       //target, AttackEvent
	OnMelt               //target, AttackEvent
	OnVaporize           //target, AttackEvent
	OnFrozen             //target, AttackEvent
	OnElectroCharged     //target, AttackEvent
	OnSwirlHydro         //target, AttackEvent
	OnSwirlCryo          //target, AttackEvent
	OnSwirlElectro       //target, AttackEvent
	OnSwirlPyro          //target, AttackEvent
	OnCrystallizeHydro   //target, AttackEvent
	OnCrystallizeCryo    //target, AttackEvent
	OnCrystallizeElectro //target, AttackEvent
	OnCrystallizePyro    //target, AttackEvent
	OnAggravate          //target, AttackEvent
	OnSpread             //target, AttackEvent
	OnQuicken            //target, AttackEvent
	OnBloom              //target, AttackEvent
	OnHyperbloom         //target, AttackEvent
	OnBurgeon            //target, AttackEvent
	OnBurning            //target, AttackEvent
	OnShatter            //target, AttackEvent; at the end to simplify all reaction event subs since it's normally not considered as an elemental reaction
	ReactionEventEndDelim
	OnDendroCore //Gadget
	//other stuff
	OnStamUse          //abil
	OnShielded         //shield
	OnCharacterSwap    //prev, next
	OnParticleReceived //particle
	OnEnergyChange     //character_received_index, pre_energy, energy_change, src (post-energy available in character_received)
	OnTargetDied       //target, AttackEvent
	OnTargetMoved      //target
	OnCharacterHit     //nil <- this is for when the character is going to get hit but might be shielded from dmg
	OnCharacterHurt    //amount
	OnHeal             //src char, target character, amount
	OnPlayerHPDrain    //DrainInfo
	OnSelfInfusion     //element, durability, duration
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
