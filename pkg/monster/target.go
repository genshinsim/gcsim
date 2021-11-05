package monster

import (
	"github.com/genshinsim/gsim/pkg/core"
)

type Target struct {
	index int
	level int
	maxHP float64
	hp    float64
	res   map[core.EleType]float64

	//modifiers
	resMod []core.ResistMod
	defMod []core.DefMod

	//icd related
	icdGroupOnTimer       [][]bool
	icdTagCounter         [][]int
	icdDamageGroupOnTimer [][]bool
	icdDamageGroupCounter [][]int

	//reactions
	aura Aura

	core *core.Core
}

func New(index int, ctrl *core.Core, p core.EnemyProfile) *Target {
	t := &Target{}

	t.index = index
	t.level = p.Level
	t.res = p.Resist
	t.core = ctrl
	t.maxHP = p.HP
	t.hp = p.HP

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

func (t *Target) AuraTick() {

}

func (t *Target) Tick() {
	//element stuff
	if t.aura != nil {
		done := t.aura.Tick()
		if done {
			if t.core.Flags.LogDebug {
				t.core.Log.Debugw(t.aura.Type().String()+" expired",
					"frame", t.core.F,
					"event", core.LogElementEvent,
					"aura", t.aura.Type(),
					"source", t.aura.Source(),
					"target", t.index,
				)
			}
			t.aura = nil
		}
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
	t.core = nil
}
