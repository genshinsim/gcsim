package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (t *Tmpl) WillApplyEle(tag core.ICDTag, grp core.ICDGroup, char int) bool {

	//no icd if no tag
	if tag == core.ICDTagNone {
		return true
	}

	//check if we need to start timer
	x := t.icdTagOnTimer[char][tag]
	if !t.icdTagOnTimer[char][tag] {
		t.icdTagOnTimer[char][tag] = true
		t.ResetTagCounterAfterDelay(tag, grp, char)
	}

	val := t.icdTagCounter[char][tag]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdTagCounter[char][tag]++
	if t.icdTagCounter[char][tag] == len(core.ICDGroupEleApplicationSequence[grp]) {
		t.icdTagCounter[char][tag] = 0
	}

	t.coretype.Log.NewEvent("ele icd check", coretype.LogICDEvent, char, "grp", grp, "target", t.TargetIndex, "tag", tag, "counter", val, "val", core.ICDGroupEleApplicationSequence[grp][val], "group on timer", x)
	//true if group seq is 1
	return core.ICDGroupEleApplicationSequence[grp][val] == 1
}

func (t *Tmpl) GroupTagDamageMult(tag core.ICDTag, grp core.ICDGroup, char int) float64 {

	//check if we need to start timer
	if !t.icdDamageTagOnTimer[char][tag] {
		t.icdDamageTagOnTimer[char][tag] = true
		t.ResetDamageCounterAfterDelay(tag, grp, char)
	}

	val := t.icdDamageTagCounter[char][tag]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdDamageTagCounter[char][tag]++
	if t.icdDamageTagCounter[char][tag] == len(core.ICDGroupDamageSequence[grp]) {
		t.icdDamageTagCounter[char][tag] = 0
	}

	//true if group seq is 1
	return core.ICDGroupDamageSequence[grp][val]
}

func (t *Tmpl) ResetDamageCounterAfterDelay(tag core.ICDTag, grp core.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdDamageTagCounter[char][tag] = 0
		t.icdDamageTagOnTimer[char][tag] = false
		t.coretype.Log.NewEvent("damage counter reset", coretype.LogICDEvent, char, "tag", tag, "grp", grp)
	}, core.ICDGroupResetTimer[grp]-1)
	t.coretype.Log.NewEvent("damage reset timer set", coretype.LogICDEvent, char, "tag", tag, "grp", grp, "reset", t.Core.Frame+core.ICDGroupResetTimer[grp]-1)
}

func (t *Tmpl) ResetTagCounterAfterDelay(tag core.ICDTag, grp core.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdTagOnTimer[char][tag] = false
		t.coretype.Log.NewEvent("ele app counter reset", coretype.LogICDEvent, char, "tag", tag, "grp", grp)
	}, core.ICDGroupResetTimer[grp]-1)
	t.coretype.Log.NewEvent("ele app reset timer set", coretype.LogICDEvent, char, "tag", tag, "grp", grp, "reset", t.Core.Frame+core.ICDGroupResetTimer[grp]-1)
}
