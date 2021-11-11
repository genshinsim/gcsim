package monster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (t *Target) handleReaction(ds *core.Snapshot) {

	// log.Println("existing aura", t.aura)

	if t.aura == nil {
		t.aura = NewAura(ds, t.core.F)
		// log.Println("new aura", t.aura)
		return
	}

	// log.Println("aura not nil, reacting")

	aura, ok := t.aura.React(ds, t)

	if ok {
		t.core.Events.Emit(core.OnReactionOccured, t, ds)
	}

	t.aura = aura

	// log.Println("result aura", t.aura, a)

}

func (t *Target) queueReaction(in *core.Snapshot, typ core.ReactionType, res core.Durability, delay int) {
	ds := t.TransReactionSnapshot(in, typ, res, true)

	t.core.Tasks.Add(func() {
		t.core.Combat.ApplyDamage(&ds)
	}, delay)

	//if swirl queue another
	switch typ {
	case core.SwirlCryo:
	case core.SwirlElectro:
	case core.SwirlHydro:
	case core.SwirlPyro:
	default:
		return
	}
	ds2 := t.TransReactionSnapshot(in, typ, res, false)
	ds2.Abil = ds2.Abil + " (aoe)"
	t.core.Tasks.Add(func() {
		t.core.Combat.ApplyDamage(&ds2)
	}, delay)
}

//calculate reaction extra damage here
func (t *Target) TransReactionSnapshot(in *core.Snapshot, typ core.ReactionType, res core.Durability, selfHarm bool) core.Snapshot {
	ds := core.Snapshot{
		CharLvl:     in.CharLvl,
		ActorEle:    in.ActorEle,
		Actor:       in.Actor,
		ActorIndex:  in.ActorIndex,
		SourceFrame: t.core.F,
		Abil:        string(typ),
		ImpulseLvl:  1,
		//reaction related
		Mult:             0,
		Durability:       0,
		StrikeType:       core.StrikeTypeDefault,
		AttackTag:        core.AttackTagNone,
		ICDTag:           core.ICDTagNone,
		ICDGroup:         core.ICDGroupReactionA,
		IsReactionDamage: true,
		ReactionType:     typ,

		//targetting info
		Targets:   core.TargetAll,
		DamageSrc: t.index,
		SelfHarm:  selfHarm,

		//no stats, should never get used
		Stats: make([]float64, core.EndStatType),
	}

	var mult float64
	switch typ {
	case core.Overload:
		mult = 2
		ds.Element = core.Pyro
		ds.AttackTag = core.AttackTagOverloadDamage
		ds.ICDTag = core.ICDTagOverloadDamage
		ds.ICDGroup = core.ICDGroupReactionB
		ds.StrikeType = core.StrikeTypeBlunt
	case core.Superconduct:
		mult = 0.5
		ds.Element = core.Cryo
		ds.AttackTag = core.AttackTagSuperconductDamage
		ds.ICDTag = core.ICDTagSuperconductDamage
		ds.OnHitCallback = superconductPhysShred
	case core.ElectroCharged:
		mult = 1.2
		ds.Element = core.Electro
		ds.AttackTag = core.AttackTagECDamage
		ds.ICDTag = core.ICDTagECDamage
		ds.ICDGroup = core.ICDGroupReactionB
	case core.Shatter:
		mult = 1.5
		ds.Element = core.Physical
	case core.SwirlElectro:
		mult = 0.6
		ds.Element = core.Electro
		ds.AttackTag = core.AttackTagSwirlElectro
		ds.ICDTag = core.ICDTagSwirlElectro
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			ds.Targets = core.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	case core.SwirlCryo:
		mult = 0.6
		ds.Element = core.Cryo
		ds.AttackTag = core.AttackTagSwirlCryo
		ds.ICDTag = core.ICDTagSwirlCryo
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			ds.Targets = core.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	case core.SwirlHydro:
		mult = 0.6
		ds.Element = core.Hydro
		ds.AttackTag = core.AttackTagSwirlHydro
		ds.ICDTag = core.ICDTagSwirlHydro
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			mult = 0
			ds.Targets = core.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	case core.SwirlPyro:
		mult = 0.6
		ds.Element = core.Pyro
		ds.AttackTag = core.AttackTagSwirlPyro
		ds.ICDTag = core.ICDTagSwirlPyro
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			ds.Targets = core.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	default:
		//either not implemented or no dmg
		return core.Snapshot{}
	}

	//grab live EM
	ds.Stats[core.EM] = t.core.Chars[in.ActorIndex].Stat(core.EM)
	ds.Mult = mult

	return ds
}

func superconductPhysShred(tar core.Target) {
	tar.AddResMod("superconductphysshred", core.ResistMod{
		Duration: 12 * 60,
		Ele:      core.Physical,
		Value:    -0.4,
	})
}

var reactionLvlBase = []float64{
	17.1656055450439,
	18.5350475311279,
	19.9048538208007,
	21.27490234375,
	22.6453990936279,
	24.6496124267578,
	26.6406421661376,
	28.8685874938964,
	31.3676795959472,
	34.1433448791503,
	37.201000213623,
	40.6599998474121,
	44.4466667175292,
	48.5635185241699,
	53.7484817504882,
	59.0818977355957,
	64.4200439453125,
	69.7244567871093,
	75.1231384277343,
	80.5847778320312,
	86.1120300292968,
	91.703742980957,
	97.24462890625,
	102.812644958496,
	108.409561157226,
	113.201690673828,
	118.102905273437,
	122.979316711425,
	129.727325439453,
	136.292907714843,
	142.670852661132,
	149.029022216796,
	155.4169921875,
	161.825500488281,
	169.106307983398,
	176.518081665039,
	184.07273864746,
	191.709518432617,
	199.556915283203,
	207.382049560546,
	215.398895263671,
	224.165664672851,
	233.502166748046,
	243.35057067871,
	256.063079833984,
	268.543487548828,
	281.526062011718,
	295.013641357421,
	309.067199707031,
	323.601593017578,
	336.757537841796,
	350.530303955078,
	364.482696533203,
	378.619171142578,
	398.600402832031,
	416.398254394531,
	434.386993408203,
	452.951049804687,
	472.606231689453,
	492.884887695312,
	513.568542480468,
	539.103210449218,
	565.510559082031,
	592.538757324218,
	624.443420410156,
	651.470153808593,
	679.496826171875,
	707.794067382812,
	736.671447753906,
	765.640258789062,
	794.773376464843,
	824.677368164062,
	851.157775878906,
	877.742065429687,
	914.229125976562,
	946.746765136718,
	979.411376953125,
	1011.22302246093,
	1044.79174804687,
	1077.44372558593,
	1109.99755859375,
	1142.9765625,
	1176.36950683593,
	1210.18444824218,
	1253.83569335937,
	1288.95275878906,
	1325.48413085937,
	1363.45690917968,
	1405.09741210937,
	1446.853515625,
}
