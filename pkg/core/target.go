package core

type Target interface {
	//basic info
	Type() TargettableType //type of target
	Index() int            //should correspond to index
	SetIndex(ind int)      //update the current index
	MaxHP() float64
	HP() float64

	//collision detection
	Shape() Shape

	//attacks
	Attack(*AttackEvent) (float64, bool)

	//reaction/aura stuff
	Tick()
	AuraContains(...EleType) bool
	AuraType() EleType

	//tags
	SetTag(key string, val int)
	GetTag(key string) int
	RemoveTag(key string)

	//target mods
	AddDefMod(key string, val float64, dur int)
	AddResMod(key string, val ResistMod)
	RemoveResMod(key string)
	RemoveDefMod(key string)
	HasDefMod(key string) bool
	HasResMod(key string) bool

	//getting rid of
	Kill()
}

type TargettableType int

const (
	TargettableEnemy TargettableType = iota
	TargettablePlayer
	TargettableObject
	TargettableTypeCount
)

// type TargetEnemy interface {
// 	Index() int
// 	SetIndex(ind int) //update the current index
// 	MaxHP() float64
// 	HP() float64
// 	//aura/reactions
// 	AuraType() EleType
// 	AuraContains(e ...EleType) bool
// 	Tick() //this should happen first before task ticks

// 	//attacks
// 	Attack(ds *Snapshot) (float64, bool)

// 	AddDefMod(key string, val float64, dur int)
// 	AddResMod(key string, val ResistMod)
// 	RemoveResMod(key string)
// 	RemoveDefMod(key string)
// 	HasDefMod(key string) bool
// 	HasResMod(key string) bool

// 	Delete() //gracefully deference everything so that it can be gc'd
// }

type ResistMod struct {
	Key      string
	Ele      EleType
	Value    float64
	Duration int
	Expiry   int
}

type DefMod struct {
	Key    string
	Value  float64
	Expiry int
}

func (c *Core) ReindexTargets() {
	//wipe out nil entries
	n := 0
	for _, v := range c.Targets {
		if v != nil {
			c.Targets[n] = v
			c.Targets[n].SetIndex(n)
			n++
		}
	}
	c.Targets = c.Targets[:n]
}
