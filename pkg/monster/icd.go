package monster

import "github.com/genshinsim/gsim/pkg/def"

func (t *Target) willApplyEle(tag def.ICDTag, grp def.ICDGroup, char int) bool {

	//no icd if no tag
	if tag == def.ICDTagNone {
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
	if t.icdTagCounter[char][tag] == len(def.ICDGroupEleApplicationSequence[grp]) {
		t.icdTagCounter[char][tag] = 0
	}

	//true if group seq is 1
	return def.ICDGroupEleApplicationSequence[grp][val] == 1
}

func (t *Target) groupTagDamageMult(grp def.ICDGroup, char int) float64 {

	//check if we need to start timer
	if !t.icdDamageGroupOnTimer[char][grp] {
		t.icdDamageGroupOnTimer[char][grp] = true
		t.resetDamageCounterAfterDelay(grp, char)
	}

	val := t.icdDamageGroupCounter[char][grp]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdDamageGroupCounter[char][grp]++
	if t.icdDamageGroupCounter[char][grp] == len(def.ICDGroupDamageSequence[grp]) {
		t.icdDamageGroupCounter[char][grp] = 0
	}

	//true if group seq is 1
	return def.ICDGroupDamageSequence[grp][val]
}

func (t *Target) resetDamageCounterAfterDelay(grp def.ICDGroup, char int) {
	t.addTask(func(t *Target) {
		//set the counter back to 0
		t.icdDamageGroupCounter[char][grp] = 0
		t.icdDamageGroupOnTimer[char][grp] = false
		t.log.Debugw("damage counter reset", "frame", t.sim.Frame(), "event", def.LogICDEvent, "char", char)
	}, def.ICDGroupResetTimer[grp]-1)
	t.log.Debugw("damage reset timer set", "frame", t.sim.Frame(), "event", def.LogICDEvent, "char", char, "grp", grp, "reset", t.sim.Frame()+def.ICDGroupResetTimer[grp]-1)
}

func (t *Target) resetTagCounterAfterDelay(tag def.ICDTag, grp def.ICDGroup, char int) {
	t.addTask(func(t *Target) {
		//set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdGroupOnTimer[char][grp] = false
		t.log.Debugw("ele app counter reset", "frame", t.sim.Frame(), "event", def.LogICDEvent, "tag", tag, "grp", grp, "char", char)
	}, def.ICDGroupResetTimer[grp]-1)
	t.log.Debugw("ele app reset timer set", "frame", t.sim.Frame(), "event", def.LogICDEvent, "tag", tag, "grp", grp, "char", char, "reset", t.sim.Frame()+def.ICDGroupResetTimer[grp]-1)
}
