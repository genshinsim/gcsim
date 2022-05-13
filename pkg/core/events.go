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
	OnParticleReceived //particle
	OnEnergyChange     //character_received, pre_energy, energy_change, src (post-energy available in character_received)
	OnTargetDied       //target
	OnCharacterHurt    //amount
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
	PreDash
	PostDash
	PreJump
	PostJump
	//sim stuff
	OnInitialize  //nil
	OnStateChange //prev, next
	OnTargetAdded //t
	EndEventTypes //elim
)
