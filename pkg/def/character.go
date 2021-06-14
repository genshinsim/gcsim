package def

type Character interface {
	Init(index int) //to be called when everything including weapon and artifacts has been loaded
	Tick()          //function to be called every frame

	//information functions
	Name() string
	CharIndex() int
	Ele() EleType
	WeaponClass() WeaponClass
	CurrentEnergy() float64 //current energy
	TalentLvlSkill() int
	TalentLvlAttack() int
	TalentLvlBurst() int
	HP() float64
	MaxHP() float64
	ModifyHP(float64)
	Stat(s StatType) float64

	//actions
	Attack(p map[string]int) int
	Aimed(p map[string]int) int
	ChargeAttack(p map[string]int) int
	HighPlungeAttack(p map[string]int) int
	LowPlungeAttack(p map[string]int) int
	Skill(p map[string]int) int
	Burst(p map[string]int) int
	Dash(p map[string]int) int

	//info methods
	ActionReady(a ActionType, p map[string]int) bool
	ActionFrames(a ActionType, p map[string]int) int
	ActionStam(a ActionType, p map[string]int) float64

	//char stat mods
	AddMod(mod CharStatMod)
	AddWeaponInfuse(inf WeaponInfusion)

	//cooldown stuff
	SetCD(a ActionType, dur int)
	Cooldown(a ActionType) int
	ResetActionCooldown(a ActionType)
	ReduceActionCooldown(a ActionType, v int)
	AddCDAdjustFunc(adj CDAdjust)

	//status stuff
	Tag(key string) int

	//energy
	ReceiveParticle(p Particle, isActive bool, partyCount int)
	AddEnergy(e float64)

	//combat
	Snapshot(name string, a AttackTag, icd ICDTag, g ICDGroup, st StrikeType, e EleType, d float64, mult float64) Snapshot
	ResetNormalCounter()
}

type CharStatMod struct {
	Key    string
	Amount func(a AttackTag) ([]float64, bool)
	Expiry int
}

type WeaponInfusion struct {
	Key    string
	Ele    EleType
	Tags   []AttackTag
	Expiry int
}

type CDAdjust struct {
	Key    string
	Amount func(a ActionType) float64
	Expiry int
}

type Particle struct {
	Source string
	Num    int
	Ele    EleType
}
