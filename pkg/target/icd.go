package target

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// icdState is the ele application and damage ICD state of one (char, tag, group) on a target.
type icdState struct {
	tagOnTimer       bool
	tagCounter       int
	damageTagOnTimer bool
	damageTagCounter int
}

func (t *Target) icdState(char int, tag attacks.ICDTag, grp attacks.ICDGroup) *icdState {
	k := NewIcdKey(char, tag, grp)
	s, ok := t.icd[k]
	if !ok {
		s = &icdState{}
		t.icd[k] = s
	}
	return s
}

func (t *Target) WillApplyEle(tag attacks.ICDTag, grp attacks.ICDGroup, char int) float64 {
	// no icd if no tag
	if tag == attacks.ICDTagNone {
		return 1
	}

	s := t.icdState(char, tag, grp)

	// check if we need to start timer
	x := s.tagOnTimer
	if !s.tagOnTimer {
		s.tagOnTimer = true
		t.ResetTagCounterAfterDelay(tag, grp, char)
	}

	val := s.tagCounter
	s.tagCounter++

	// if counter > length, then use 0 for group seq
	groupSeq := attacks.ICDGroupEleApplicationSequence[grp][len(attacks.ICDGroupEleApplicationSequence[grp])-1]
	if val < len(attacks.ICDGroupEleApplicationSequence[grp]) {
		groupSeq = attacks.ICDGroupEleApplicationSequence[grp][val]
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
	s := t.icdState(char, tag, grp)

	// check if we need to start timer
	if !s.damageTagOnTimer {
		s.damageTagOnTimer = true
		t.ResetDamageCounterAfterDelay(tag, grp, char)
	}

	val := s.damageTagCounter
	s.damageTagCounter++

	// if counter > length, then use 0 for group seq
	groupSeq := attacks.ICDGroupDamageSequence[grp][len(attacks.ICDGroupDamageSequence[grp])-1]
	if val < len(attacks.ICDGroupDamageSequence[grp]) {
		groupSeq = attacks.ICDGroupDamageSequence[grp][val]
	}

	return groupSeq
}

func (t *Target) ResetDamageCounterAfterDelay(tag attacks.ICDTag, grp attacks.ICDGroup, char int) {
	s := t.icdState(char, tag, grp)
	t.Core.Tasks.Add(func() {
		// set the counter back to 0
		s.damageTagCounter = 0
		s.damageTagOnTimer = false
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
	s := t.icdState(char, tag, grp)
	t.Core.Tasks.Add(func() {
		// set the counter back to 0
		s.tagCounter = 0
		s.tagOnTimer = false
		t.Core.Log.NewEvent("ele app counter reset", glog.LogICDEvent, char).
			Write("tag", tag).
			Write("grp", grp)
	}, attacks.ICDGroupResetTimer[grp]-1)
	t.Core.Log.NewEvent("ele app reset timer set", glog.LogICDEvent, char).
		Write("tag", tag).
		Write("grp", grp).
		Write("reset", t.Core.F+attacks.ICDGroupResetTimer[grp]-1)
}
