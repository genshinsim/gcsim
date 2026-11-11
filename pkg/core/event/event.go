package event

type Event int

const (
	OnEnemyHit     Event = iota // target, AttackEvent
	OnPlayerHit                 // char, AttackEvent; emits before the player/shields takes damage
	OnGadgetHit                 // target, AttackEvent
	OnEnemyDamage               // target, AttackEvent, amount, crit
	OnGadgetDamage              // target, AttackEvent
	OnApplyAttack               // AttackEvent
	// reaction related
	// OnReactionOccured // target, AttackEvent
	// OnTransReaction   // target, AttackEvent
	// OnAmpReaction     // target, AttackEvent
	OnElementApplied         // target, AttackEvent
	OnAuraDurabilityAdded    // target, ele, durability
	OnAuraDurabilityDepleted // target, ele
	// OnReaction               // target, AttackEvent, ReactionType
	ReactionEventStartDelim
	OnOverload           // target, AttackEvent
	OnSuperconduct       // target, AttackEvent
	OnMelt               // target, AttackEvent
	OnVaporize           // target, AttackEvent
	OnFrozen             // target, AttackEvent
	OnElectroCharged     // target, AttackEvent
	OnSwirlHydro         // target, AttackEvent
	OnSwirlCryo          // target, AttackEvent
	OnSwirlElectro       // target, AttackEvent
	OnSwirlPyro          // target, AttackEvent
	OnCrystallizeHydro   // target, AttackEvent
	OnCrystallizeCryo    // target, AttackEvent
	OnCrystallizeElectro // target, AttackEvent
	OnCrystallizePyro    // target, AttackEvent
	OnAggravate          // target, AttackEvent
	OnSpread             // target, AttackEvent
	OnQuicken            // target, AttackEvent
	OnBloom              // target, AttackEvent
	OnHyperbloom         // target, AttackEvent
	OnBurgeon            // target, AttackEvent
	OnBurning            // target, AttackEvent
	OnLunarCharged       // target, AttackEvent
	OnLunarBloom         // target, AttackEvent
	OnLunarCrystallize   // target, AttackEvent
	OnStellarConduct     // target, AttackEvent
	OnStellarSwirl       // target, AttackEvent
	OnShatter            // target, AttackEvent; at the end to simplify all reaction event subs since it's normally not considered as an elemental reaction
	ReactionEventEndDelim
	OnDendroCore // Gadget
	// other stuff
	OnStamUse           // abil
	OnShielded          // shield
	OnShieldBreak       // shield break
	OnConstructSpawned  // nil
	OnCharacterSwap     // prev, next
	OnParticleReceived  // particle
	OnEnergyChange      // character_received, pre_energy, energy_change, src (post-energy available in character_received), is_particle (boolean)
	OnEnergyBurst       // character_drained, pre_energy, burst_cost
	OnTargetDied        // target, AttackEvent
	OnTargetMoved       // target
	OnCharacterHurt     // amount
	OnHPDebt            // target character, amount
	OnHeal              // src char, target character, amount, overheal, amount_before_debt
	OnPlayerPreHPDrain  // Draininfo to modify
	OnPlayerHPDrain     // DrainInfo
	OnNightsoulBurst    // target, AttackEvent
	OnNightsoulGenerate // char, amount
	OnNightsoulConsume  // char, amount
	OnMovement          // float64; distance moved
	// ability use
	OnActionFailed // ActiveCharIndex, action.Action, param, action.ActionFailure
	OnActionExec   // ActiveCharIndex, action.Action, param
	OnSkill        // nil
	OnBurst        // nil
	OnAttack       // nil
	OnChargeAttack // nil
	OnPlunge       // nil
	OnAimShoot     // nil
	OnDash
	OnSpecialReactionAttack // target, AttackEvent; event so predamagemods can be applied to the individual lunar/stellar contributions. Emitted once per contributor
	OnMoondriftHarmony      // target, AttackEvent;
	OnStellarVortexDetonate // src char, contribMap, AttackPattern;
	// sim stuff
	OnInitialize  // nil
	OnStateChange // prev, next, segmented
	OnEnemyAdded  // t
	OnTick
	OnSimEndedSuccessfully // nil
	EndEventTypes          // elim
)

type Handler struct {
	events [][]ehook
	// argBufs[d] holds the args handed to hooks by an Emit nested d deep. Reusing them
	// means Emit doesn't let args escape, so callers keep the variadic slice on the stack
	// instead of allocating one per emit. Hooks must not keep args after returning.
	argBufs [][]any
	depth   int
}

// Hook handles an emitted event. args is a buffer Emit reuses for the next event, so a hook must
// not keep args, or a pointer into it, after returning: copy out the elements it needs instead.
// Assigning to an element doesn't reach the emitter, so hooks don't; modifying what an element
// points to is fine. TestHookArgs enforces this.
type Hook func(args ...any)

type Eventter interface {
	Subscribe(e Event, f Hook, key string)
	Unsubscribe(e Event, key string)
	Emit(e Event, args ...any)
}

type ehook struct {
	f   Hook
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

// Subscribe to an event
//
//	core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
//		e, ok := args[0].(*enemy.Enemy);
//		if !ok {
//			return
//		}
//		atk, ok := args[1].(*info.AttackEvent)
//		if !ok {
//			return
//		}
//		if atk.Info.ActorIndex != char.Index() {
//			return
//		}
//		doSomething(e)
//	}, "subscription")
func (h *Handler) Subscribe(e Event, f Hook, key string) {
	a := h.events[e]

	evt := ehook{
		f:   f,
		key: key,
	}

	// check if override first
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
	for i, v := range h.events[e] {
		if v.key == key {
			h.events[e][i].f = nil
		}
	}
}

// HasHooks reports whether anything subscribed to e. Without hooks Emit(e, ...) does nothing.
func (h *Handler) HasHooks(e Event) bool {
	return len(h.events[e]) > 0
}

func (h *Handler) Emit(e Event, args ...any) {
	hooks := h.events[e]
	if len(hooks) == 0 {
		return
	}
	// a hook may emit again, so each nesting depth gets its own buffer
	if h.depth == len(h.argBufs) {
		h.argBufs = append(h.argBufs, make([]any, 7))
	}
	buf := h.argBufs[h.depth]
	if len(args) > len(buf) {
		buf = make([]any, len(args))
		h.argBufs[h.depth] = buf
	}
	// cap == len so a hook's append(args, x) copies instead of writing into the spare capacity
	buf = buf[:len(args):len(args)]
	// element-wise rather than copy/clear: those are runtime calls on wasm, and args is short
	for i, a := range args { //nolint:staticcheck // S1001: see above
		buf[i] = a
	}
	h.depth++
	for _, v := range hooks {
		if v.f != nil {
			v.f(buf...)
		}
	}
	h.depth--
	for i := 0; i < len(buf); i++ { // not `range`, which compiles to a memclr call
		buf[i] = nil
	}
}
