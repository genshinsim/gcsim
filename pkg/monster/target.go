package monster

import (
	"math/rand"

	"github.com/genshinsim/gsim/pkg/def"
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

	//damage related
	onAttackLandedFuncs []attackLandedFunc

	//reactions
	aura              Aura
	onReactionOccured []reactionHooks //reaction hooks

	sim  def.Sim
	rand *rand.Rand
	log  def.Logger
}

func New(index int, s def.Sim, log def.Logger, hp float64, p def.EnemyProfile) *Target {
	t := &Target{}

	t.index = index
	t.level = p.Level
	t.res = p.Resist
	t.log = log
	t.sim = s
	t.rand = s.Rand()
	t.maxHP = hp
	t.hp = hp

	t.icdGroupOnTimer = make([][]bool, def.MaxTeamPlayerCount)
	t.icdTagCounter = make([][]int, def.MaxTeamPlayerCount)
	t.icdDamageGroupCounter = make([][]int, def.MaxTeamPlayerCount)
	t.icdDamageGroupOnTimer = make([][]bool, def.MaxTeamPlayerCount)

	for i := 0; i < 4; i++ {
		t.icdGroupOnTimer[i] = make([]bool, def.ICDGroupLength)
		t.icdTagCounter[i] = make([]int, def.ICDTagLength)
		t.icdDamageGroupCounter[i] = make([]int, def.ICDGroupLength)
		t.icdDamageGroupOnTimer[i] = make([]bool, def.ICDGroupLength)
	}

	t.tasks = make(map[int][]func(t *Target))

	t.onAttackLandedFuncs = make([]attackLandedFunc, 0, 10)

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

func (t *Target) AuraTick() {
	//element stuff
	if t.aura != nil {
		done := t.aura.Tick()
		if done {
			t.aura = nil
		}
	}
}

func (t *Target) Tick() {
	//check tasks
	for _, f := range t.tasks[t.sim.Frame()] {
		f(t)
	}

	delete(t.tasks, t.sim.Frame())

}

type attackLandedFunc struct {
	f func(ds *def.Snapshot)
	k string
}

func (t *Target) AddOnAttackLandedHook(fun func(ds *def.Snapshot), key string) {
	ind := -1
	for i, v := range t.onAttackLandedFuncs {
		if v.k == key {
			ind = i
		}
	}
	if ind != -1 {
		t.onAttackLandedFuncs[ind] = attackLandedFunc{
			f: fun,
			k: key,
		}
		return
	}
	t.onAttackLandedFuncs = append(t.onAttackLandedFuncs, attackLandedFunc{
		f: fun,
		k: key,
	})
}

func (t *Target) RemoveOnAttackLandedHook(key string) {
	ind := -1
	for i, v := range t.onAttackLandedFuncs {
		if v.k == key {
			ind = i
		}
	}
	if ind != -1 {
		t.onAttackLandedFuncs[ind] = t.onAttackLandedFuncs[len(t.onAttackLandedFuncs)-1]
		t.onAttackLandedFuncs = t.onAttackLandedFuncs[:len(t.onAttackLandedFuncs)-1]
	}
}

func (t *Target) onAttackLanded(ds *def.Snapshot) {
	for _, v := range t.onAttackLandedFuncs {
		v.f(ds)
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
