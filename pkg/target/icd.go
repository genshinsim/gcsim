package target

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (t *Target) WillApplyEle(tag attacks.ICDTag, grp attacks.ICDGroup, char int) float64 {
	// no icd if no tag
	if tag == attacks.ICDTagNone {
		return 1
	}

	// check if we need to start timer
	x := t.icdTagOnTimer[char][tag]
	if !t.icdTagOnTimer[char][tag] {
		t.icdTagOnTimer[char][tag] = true
		t.ResetTagCounterAfterDelay(tag, grp, char)
	}

	val := t.icdTagCounter[char][tag]
	t.icdTagCounter[char][tag]++

	// if counter > length, then use 0 for group seq
	groupSeq := combat.ICDGroupEleApplicationSequence[grp][len(combat.ICDGroupEleApplicationSequence[grp])-1]
	if val < len(combat.ICDGroupEleApplicationSequence[grp]) {
		groupSeq = combat.ICDGroupEleApplicationSequence[grp][val]
	}

	t.Core.Log.NewEvent("ele icd check", glog.LogICDEvent, char).
		Write("grp", grp).
		Write("target", t.key).
		Write("tag", tag).
		Write("counter", val).
		Write("val", groupSeq).
		Write("group on timer", x)

	return groupSeq
}

func (t *Target) GroupTagDamageMult(tag attacks.ICDTag, grp attacks.ICDGroup, char int) float64 {
	// check if we need to start timer
	if !t.icdDamageTagOnTimer[char][tag] {
		t.icdDamageTagOnTimer[char][tag] = true
		t.ResetDamageCounterAfterDelay(tag, grp, char)
	}

	val := t.icdDamageTagCounter[char][tag]
	t.icdDamageTagCounter[char][tag]++

	// if counter > length, then use 0 for group seq
	groupSeq := combat.ICDGroupDamageSequence[grp][len(combat.ICDGroupDamageSequence[grp])-1]
	if val < len(combat.ICDGroupDamageSequence[grp]) {
		groupSeq = combat.ICDGroupDamageSequence[grp][val]
	}

	return groupSeq
}

func (t *Target) ResetDamageCounterAfterDelay(tag attacks.ICDTag, grp attacks.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		// set the counter back to 0
		t.icdDamageTagCounter[char][tag] = 0
		t.icdDamageTagOnTimer[char][tag] = false
		t.Core.Log.NewEvent("damage counter reset", glog.LogICDEvent, char).
			Write("tag", tag).
			Write("grp", grp)
	}, attacks.ICDGroupResetTimer[grp]-1)
	t.Core.Log.NewEvent("damage reset timer set", glog.LogICDEvent, char).
		Write("tag", tag).
		Write("grp", grp).
		Write("reset", t.Core.F+attacks.ICDGroupResetTimer[grp]-1)
}

func (t *Target) ResetTagCounterAfterDelay(tag attacks.ICDTag, grp attacks.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		// set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdTagOnTimer[char][tag] = false
		t.Core.Log.NewEvent("ele app counter reset", glog.LogICDEvent, char).
			Write("tag", tag).
			Write("grp", grp)
	}, attacks.ICDGroupResetTimer[grp]-1)
	t.Core.Log.NewEvent("ele app reset timer set", glog.LogICDEvent, char).
		Write("tag", tag).
		Write("grp", grp).
		Write("reset", t.Core.F+attacks.ICDGroupResetTimer[grp]-1)
}
