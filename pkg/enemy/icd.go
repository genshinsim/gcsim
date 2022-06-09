package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (t *Enemy) WillApplyEle(tag combat.ICDTag, grp combat.ICDGroup, char int) bool {

	//no icd if no tag
	if tag == combat.ICDTagNone {
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
	if t.icdTagCounter[char][tag] == len(combat.ICDGroupEleApplicationSequence[grp]) {
		t.icdTagCounter[char][tag] = 0
	}

	t.Core.Log.NewEvent("ele icd check", glog.LogICDEvent, char, "grp", grp, "target", t.TargetIndex, "tag", tag, "counter", val, "val", combat.ICDGroupEleApplicationSequence[grp][val], "group on timer", x)
	//true if group seq is 1
	return combat.ICDGroupEleApplicationSequence[grp][val] == 1
}

func (t *Enemy) GroupTagDamageMult(tag combat.ICDTag, grp combat.ICDGroup, char int) float64 {

	//check if we need to start timer
	if !t.icdDamageTagOnTimer[char][tag] {
		t.icdDamageTagOnTimer[char][tag] = true
		t.ResetDamageCounterAfterDelay(tag, grp, char)
	}

	val := t.icdDamageTagCounter[char][tag]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdDamageTagCounter[char][tag]++
	if t.icdDamageTagCounter[char][tag] == len(combat.ICDGroupDamageSequence[grp]) {
		t.icdDamageTagCounter[char][tag] = 0
	}

	//true if group seq is 1
	return combat.ICDGroupDamageSequence[grp][val]
}

func (t *Enemy) ResetDamageCounterAfterDelay(tag combat.ICDTag, grp combat.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdDamageTagCounter[char][tag] = 0
		t.icdDamageTagOnTimer[char][tag] = false
		t.Core.Log.NewEvent("damage counter reset", glog.LogICDEvent, char, "tag", tag, "grp", grp)
	}, combat.ICDGroupResetTimer[grp]-1)
	t.Core.Log.NewEvent("damage reset timer set", glog.LogICDEvent, char, "tag", tag, "grp", grp, "reset", t.Core.F+combat.ICDGroupResetTimer[grp]-1)
}

func (t *Enemy) ResetTagCounterAfterDelay(tag combat.ICDTag, grp combat.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdTagOnTimer[char][tag] = false
		t.Core.Log.NewEvent("ele app counter reset", glog.LogICDEvent, char, "tag", tag, "grp", grp)
	}, combat.ICDGroupResetTimer[grp]-1)
	t.Core.Log.NewEvent("ele app reset timer set", glog.LogICDEvent, char, "tag", tag, "grp", grp, "reset", t.Core.F+combat.ICDGroupResetTimer[grp]-1)
}
