package monster

import "github.com/genshinsim/gsim/pkg/core"

func (t *Target) willApplyEle(tag core.ICDTag, grp core.ICDGroup, char int) bool {

	//no icd if no tag
	if tag == core.ICDTagNone {
		return true
	}

	//check if we need to start timer
	if !t.icdGroupOnTimer[char][grp] {
		t.icdGroupOnTimer[char][grp] = true
		t.resetTagCounterAfterDelay(tag, grp, char)
	}

	val := t.icdTagCounter[char][tag]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdTagCounter[char][tag]++
	if t.icdTagCounter[char][tag] == len(core.ICDGroupEleApplicationSequence[grp]) {
		t.icdTagCounter[char][tag] = 0
	}

	//true if group seq is 1
	if core.ICDGroupEleApplicationSequence[grp][val] != 1 {
		t.core.Log.Debugw("ele app on icd", "frame", t.core.F, "event", core.LogICDEvent, "char", char, "target", t.index)
		return false
	}

	return true
}

func (t *Target) groupTagDamageMult(grp core.ICDGroup, char int) float64 {

	//check if we need to start timer
	if !t.icdDamageGroupOnTimer[char][grp] {
		t.icdDamageGroupOnTimer[char][grp] = true
		t.resetDamageCounterAfterDelay(grp, char)
	}

	val := t.icdDamageGroupCounter[char][grp]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdDamageGroupCounter[char][grp]++
	if t.icdDamageGroupCounter[char][grp] == len(core.ICDGroupDamageSequence[grp]) {
		t.icdDamageGroupCounter[char][grp] = 0
	}

	//true if group seq is 1
	if core.ICDGroupDamageSequence[grp][val] == 0 {
		t.core.Log.Debugw("dmg on icd", "frame", t.core.F, "event", core.LogICDEvent, "char", char, "target", t.index)
		return 0
	}
	return core.ICDGroupDamageSequence[grp][val]
}

func (t *Target) resetDamageCounterAfterDelay(grp core.ICDGroup, char int) {
	t.core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdDamageGroupCounter[char][grp] = 0
		t.icdDamageGroupOnTimer[char][grp] = false
		t.core.Log.Debugw("damage counter reset", "frame", t.core.F, "event", core.LogICDEvent, "char", char)
	}, core.ICDGroupResetTimer[grp]-1)
	t.core.Log.Debugw("damage reset timer set", "frame", t.core.F, "event", core.LogICDEvent, "char", char, "grp", grp, "reset", t.core.F+core.ICDGroupResetTimer[grp]-1)
}

func (t *Target) resetTagCounterAfterDelay(tag core.ICDTag, grp core.ICDGroup, char int) {
	t.core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdGroupOnTimer[char][grp] = false
		t.core.Log.Debugw("ele app counter reset", "frame", t.core.F, "event", core.LogICDEvent, "tag", tag, "grp", grp, "char", char)
	}, core.ICDGroupResetTimer[grp]-1)
	t.core.Log.Debugw("ele app reset timer set", "frame", t.core.F, "event", core.LogICDEvent, "tag", tag, "grp", grp, "char", char, "reset", t.core.F+core.ICDGroupResetTimer[grp]-1)
}
