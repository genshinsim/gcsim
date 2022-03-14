package coretype

type Framer interface {
	F() int
}

type CommandExecuter interface {
	Exec(n Command) (frames int, done bool, err error) //return frames, if executed, any errors
}

type Logger interface {
	NewEvent(msg string, typ LogSource, srcChar int, keysAndValues ...interface{}) LogEvent
	NewEventBuildMsg(typ LogSource, srcChar int, msg ...string) LogEvent
	Dump() ([]byte, error) //print out all the logged events in array of JSON strings in the ordered they were added
}

type Statuser interface {
	StatusDuration(key string) int
	AddStatus(key string, dur int)
	ExtendStatus(key string, dur int)
	DeleteStatus(key string)
}

type EventEmitter interface {
	Subscribe(e EventType, f EventHook, key string)
	Unsubscribe(e EventType, key string)
	Emit(e EventType, args ...interface{})
}

type EventHook func(args ...interface{}) bool

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
	AddTask(f func(), delay int)
}

type ConstructHandler interface {
	NewConstruct(c Construct, refresh bool)
	NewConstructNoLimit(c Construct, refresh bool)
	CountConstruct() int
	CountConstructByType(t GeoConstructType) int
	DestroyConstruct(key int) bool
	HasConstruct(key int) bool
	ConstructExpiry(t GeoConstructType) int
}

type QueueHandler interface {
	//returns a sequence of 1 or more commands to execute,
	//whether or not to drop sequence if any is not ready, and any error
	Next() (queue []Command, dropIfFailed bool, err error)
	SetActionList(pq []ActionBlock) error
}

// The rest are for players

type EnergyHandler interface {
	DistributeParticle(p Particle)
}
type HealthHandler interface {
	Heal(hi HealInfo)
	AddIncHealBonus(f func(healedCharIndex int) float64)

	AddDamageReduction(f func() (float64, bool))
	HurtChar(dmg float64, ele EleType)
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
