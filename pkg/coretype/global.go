package coretype

type Tagger interface {
	SetTag(key string, val int)
	GetTag(key string) int
	RemoveTag(key string)
}

type Ticker interface {
	Tick()
}

type Initer interface {
	Init()
}

type Indexer interface {
	Index() int
	SetIndex(ind int)
}

type Reactable interface {
	//reaction/aura stuff
	Tick()
	AuraContains(...EleType) bool
	AuraType() EleType
}

type Target interface {
	//basic info
	Shape() Shape //collision detection
	MaxHP() float64
	HP() float64

	//attacks
	Attack(*AttackEvent, LogEvent) (float64, bool)

	//getting rid of
	Kill()
}
