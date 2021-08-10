package monster

import (
	"math/rand"

	"github.com/genshinsim/gsim/pkg/core"
)

type Target struct {
	index int
	level int
	maxHP float64
	hp    float64
	res   map[core.EleType]float64
	tasks map[int][]func(t *Target)

	//modifiers
	resMod []core.ResistMod
	defMod []core.DefMod

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

	sim  core.Sim
	rand *rand.Rand
	log  core.Logger
}

func New(index int, s core.Sim, log core.Logger, hp float64, p core.EnemyProfile) *Target {
	t := &Target{}

	t.index = index
	t.level = p.Level
	t.res = p.Resist
	t.log = log
	t.sim = s
	t.rand = s.Rand()
	t.maxHP = hp
	t.hp = hp

	t.icdGroupOnTimer = make([][]bool, core.MaxTeamPlayerCount)
	t.icdTagCounter = make([][]int, core.MaxTeamPlayerCount)
	t.icdDamageGroupCounter = make([][]int, core.MaxTeamPlayerCount)
	t.icdDamageGroupOnTimer = make([][]bool, core.MaxTeamPlayerCount)

	for i := 0; i < 4; i++ {
		t.icdGroupOnTimer[i] = make([]bool, core.ICDGroupLength)
		t.icdTagCounter[i] = make([]int, core.ICDTagLength)
		t.icdDamageGroupCounter[i] = make([]int, core.ICDGroupLength)
		t.icdDamageGroupOnTimer[i] = make([]bool, core.ICDGroupLength)
	}

	t.tasks = make(map[int][]func(t *Target))

	t.onAttackLandedFuncs = make([]attackLandedFunc, 0, 10)

	return t
}

func (t *Target) AuraType() core.EleType {
	if t.aura == nil {
		return core.NoElement
	}
	return t.aura.Type()
}

func (t *Target) AuraContains(ele ...core.EleType) bool {
	if t.aura == nil {
		for _, v := range ele {
			if v == core.NoElement {
				return true
			}
		}
		return false
	}
	return t.aura.AuraContains(ele...)
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
	f func(ds *core.Snapshot)
	k string
}

func (t *Target) AddOnAttackLandedHook(fun func(ds *core.Snapshot), key string) {
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

func (t *Target) onAttackLanded(ds *core.Snapshot) {
	for _, v := range t.onAttackLandedFuncs {
		v.f(ds)
	}
}

func (t *Target) Index() int {
	return t.index
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
