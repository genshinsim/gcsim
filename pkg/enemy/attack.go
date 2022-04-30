package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (t *Enemy) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	//if target is frozen prior to attack landing, set impulse to 0
	//let the break freeze attack to trigger actual impulse
	if t.AuraType() == attributes.Frozen {
		atk.Info.NoImpulse = true
	}

	//check shatter first
	t.ShatterCheck(atk)

	//check tags
	if atk.Info.Durability > 0 {
		//check for ICD first
		if t.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex) && atk.Info.Element != attributes.Physical {
			existing := t.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			t.React(atk)
			if t.Core.Flags.LogDebug {
				t.Core.Log.NewEvent(
					"application",
					glog.LogElementEvent,
					atk.Info.ActorIndex,
					"attack_tag", atk.Info.AttackTag,
					"applied_ele", atk.Info.Element.String(),
					"dur", applied,
					"abil", atk.Info.Abil,
					"target", t.TargetIndex,
					"existing", existing,
					"after", t.Reactable.ActiveAuraString(),
				)
			}
		}
	}

	damage, isCrit := t.calc(atk, evt)

	//record dmg
	t.HPCurrent -= damage

	return damage, isCrit
}
