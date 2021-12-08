package target

import "github.com/genshinsim/gcsim/pkg/core"

func (t *Tmpl) WillApplyEle(tag core.ICDTag, grp core.ICDGroup, char int) bool {

	//no icd if no tag
	if tag == core.ICDTagNone {
		return true
	}

	//check if we need to start timer
	if !t.icdGroupOnTimer[char][grp] {
		t.icdGroupOnTimer[char][grp] = true
		t.ResetTagCounterAfterDelay(tag, grp, char)
	}

	val := t.icdTagCounter[char][tag]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdTagCounter[char][tag]++
	if t.icdTagCounter[char][tag] == len(core.ICDGroupEleApplicationSequence[grp]) {
		t.icdTagCounter[char][tag] = 0
	}

	//true if group seq is 1
	return core.ICDGroupEleApplicationSequence[grp][val] == 1
}

func (t *Tmpl) GroupTagDamageMult(grp core.ICDGroup, char int) float64 {

	//check if we need to start timer
	if !t.icdDamageGroupOnTimer[char][grp] {
		t.icdDamageGroupOnTimer[char][grp] = true
		t.ResetDamageCounterAfterDelay(grp, char)
	}

	val := t.icdDamageGroupCounter[char][grp]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdDamageGroupCounter[char][grp]++
	if t.icdDamageGroupCounter[char][grp] == len(core.ICDGroupDamageSequence[grp]) {
		t.icdDamageGroupCounter[char][grp] = 0
	}

	//true if group seq is 1
	return core.ICDGroupDamageSequence[grp][val]
}

func (t *Tmpl) ResetDamageCounterAfterDelay(grp core.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdDamageGroupCounter[char][grp] = 0
		t.icdDamageGroupOnTimer[char][grp] = false
		t.Core.Log.Debugw("damage counter reset", "frame", t.Core.F, "event", core.LogICDEvent, "char", char)
	}, core.ICDGroupResetTimer[grp]-1)
	t.Core.Log.Debugw("damage reset timer set", "frame", t.Core.F, "event", core.LogICDEvent, "char", char, "grp", grp, "reset", t.Core.F+core.ICDGroupResetTimer[grp]-1)
}

func (t *Tmpl) ResetTagCounterAfterDelay(tag core.ICDTag, grp core.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdGroupOnTimer[char][grp] = false
		t.Core.Log.Debugw("ele app counter reset", "frame", t.Core.F, "event", core.LogICDEvent, "tag", tag, "grp", grp, "char", char)
	}, core.ICDGroupResetTimer[grp]-1)
	t.Core.Log.Debugw("ele app reset timer set", "frame", t.Core.F, "event", core.LogICDEvent, "tag", tag, "grp", grp, "char", char, "reset", t.Core.F+core.ICDGroupResetTimer[grp]-1)
}
