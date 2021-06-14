package monster

import (
	"math/rand"

	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

type Target struct {
	index int
	level int
	maxHP float64
	hp    float64
	res   map[def.EleType]float64
	tasks map[int][]func(t *Target)

	//modifiers
	resMod []def.ResistMod
	defMod []def.DefMod

	//icd related
	icdGroupOnTimer       [][]bool
	icdTagCounter         [][]int
	icdDamageGroupOnTimer [][]bool
	icdDamageGroupCounter [][]int

	//reactions
	auras             []Aura
	onReactionOccured []reactionHooks //reaction hooks

	sim  def.Sim
	rand *rand.Rand
	log  *zap.SugaredLogger
}

func New(index int, s def.Sim, log *zap.SugaredLogger, p def.EnemyProfile) *Target {
	t := &Target{}

	t.index = index
	t.level = p.Level
	t.res = p.Resist
	t.log = log
	t.sim = s
	t.rand = s.Rand()

	t.auras = make([]Aura, def.ElementMaxCount)

	for i := 0; i < 4; i++ {
		t.icdGroupOnTimer[i] = make([]bool, def.ICDGroupLength)
		t.icdTagCounter[i] = make([]int, def.ICDTagLength)
		t.icdDamageGroupCounter[i] = make([]int, def.ICDGroupLength)
		t.icdDamageGroupOnTimer[i] = make([]bool, def.ICDGroupLength)
	}

	return t
}

func (t *Target) HP() float64 {
	return t.hp
}

func (t *Target) MaxHP() float64 {
	return t.maxHP
}

func (t *Target) addTask(fun func(t *Target), delay int) {
	f := t.sim.Frame()
	t.tasks[f+delay] = append(t.tasks[f+delay], fun)
}

func (t *Target) Tick() {
	//check tasks
	for _, f := range t.tasks[t.sim.Frame()] {
		f(t)
	}

	delete(t.tasks, t.sim.Frame())

	//element stuff
	for _, a := range t.auras {
		if a != nil {
			a.Tick()
		}
	}
}

func (t *Target) SetIndex(ind int) {
	t.index = ind
}

func (t *Target) Delete() {
	t.res = nil
	t.resMod = nil
	t.defMod = nil
	t.sim = nil
	t.rand = nil
	t.log = nil
}
