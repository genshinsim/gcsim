package core

type CommandHandler interface {
	Exec(n Command) (frames int, done bool, err error) //return frames, if executed, any errors
}

type StatusHandler interface {
	Duration(key string) int
	AddStatus(key string, dur int)
	ExtendStatus(key string, dur int)
	DeleteStatus(key string)
}

type EventHandler interface {
	Subscribe(e EventType, f EventHook, key string)
	Unsubscribe(e EventType, key string)
	Emit(e EventType, args ...interface{})
}

type EventHook func(args ...interface{}) bool

type EnergyHandler interface {
	DistributeParticle(p Particle)
}

type CombatHandler interface {
	ApplyDamage(*AttackEvent) float64
	QueueAttack(a AttackInfo, p AttackPattern, snapshotDelay int, dmgDelay int, callbacks ...AttackCBFunc)
	QueueAttackWithSnap(a AttackInfo, s Snapshot, p AttackPattern, dmgDelay int, callbacks ...AttackCBFunc)
	QueueAttackEvent(ae *AttackEvent, dmgDelay int)
	TargetHasResMod(debuff string, param int) bool
	TargetHasDefMod(debuff string, param int) bool
	TargetHasElement(ele EleType, param int) bool
}

type TaskHandler interface {
	Add(f func(), delay int)
	Run()
}

type ConstructHandler interface {
	New(c Construct, refresh bool)
	NewNoLimitCons(c Construct, refresh bool)
	Count() int
	CountByType(t GeoConstructType) int
	Destroy(key int) bool
	Has(key int) bool
	Expiry(t GeoConstructType) int
	Tick()
}

type HealthHandler interface {
	Heal(hi HealInfo)
	Drain(di DrainInfo)
	AddIncHealBonus(f func(healedCharIndex int) float64)

	AddDamageReduction(f func() (float64, bool))
	HurtChar(dmg float64, ele EleType)
}

type QueueHandler interface {
	//returns a sequence of 1 or more commands to execute,
	//whether or not to drop sequence if any is not ready, and any error
	Next() (queue []Command, dropIfFailed bool, err error)
	SetActionList(pq []ActionBlock) error
}

type ShieldHandler interface {
	Add(shd Shield)
	IsShielded(char int) bool
	Get(t ShieldType) Shield
	AddBonus(f func() float64)
	OnDamage(dmg float64, ele EleType) float64
	Count() int
	Tick()
}
