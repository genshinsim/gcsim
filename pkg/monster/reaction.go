package monster

import (
	"github.com/genshinsim/gsim/pkg/def"
)

func (t *Target) handleReaction(ds *def.Snapshot) {

	// log.Println("existing aura", t.aura)

	if t.aura == nil {
		t.aura = NewAura(ds, t.sim.Frame())
		// log.Println("new aura", t.aura)
		return
	}

	// log.Println("aura not nil, reacting")

	aura, ok := t.aura.React(ds, t)

	if ok {
		t.sim.OnReaction(t, ds)
		for _, v := range t.onReactionOccured {
			v.f(ds)
		}
	}

	t.aura = aura

	// log.Println("result aura", t.aura, a)

}

type reactionHooks struct {
	f   func(ds *def.Snapshot)
	key string
}

func (t *Target) AddOnReactionHook(key string, fun func(ds *def.Snapshot)) {
	a := t.onReactionOccured

	ind := len(a)

	for i, v := range a {
		if v.key == key {
			ind = i
			break
		}
	}
	if ind != 0 && ind != len(a) {
		a[ind] = reactionHooks{
			key: key,
			f:   fun,
		}
		t.log.Debugw("reaction hook added", "frame", t.sim.Frame(), "event", def.LogHookEvent, "overwrite", true, "key", key, "hooks", a)
	} else {
		a = append(a, reactionHooks{
			key: key,
			f:   fun,
		})
		t.log.Debugw("reaction hook added", "frame", t.sim.Frame(), "event", def.LogHookEvent, "overwrite", true, "key", key, "hooks", a)
	}
	t.onReactionOccured = a
}

func (t *Target) queueReaction(in *def.Snapshot, typ def.ReactionType, res def.Durability, delay int) {
	ds := t.TransReactionSnapshot(in, typ, res, true)

	t.addTask(func(t *Target) {
		t.sim.ApplyDamage(&ds)
	}, delay)
	//if swirl queue another
	switch typ {
	case def.SwirlCryo:
	case def.SwirlElectro:
	case def.SwirlHydro:
	case def.SwirlPyro:
	default:
		return
	}
	ds2 := t.TransReactionSnapshot(in, typ, res, false)
	t.addTask(func(t *Target) {
		t.sim.ApplyDamage(&ds2)
	}, delay)
}

//calculate reaction extra damage here
func (t *Target) TransReactionSnapshot(in *def.Snapshot, typ def.ReactionType, res def.Durability, selfHarm bool) def.Snapshot {
	ds := def.Snapshot{
		CharLvl:     in.CharLvl,
		ActorEle:    in.ActorEle,
		Actor:       in.Actor,
		ActorIndex:  in.ActorIndex,
		SourceFrame: t.sim.Frame(),
		Abil:        string(typ),
		ImpulseLvl:  1,
		//reaction related
		Mult:             0,
		Durability:       0,
		StrikeType:       def.StrikeTypeDefault,
		AttackTag:        def.AttackTagNone,
		ICDTag:           def.ICDTagNone,
		ICDGroup:         def.ICDGroupReactionA,
		IsReactionDamage: true,

		//targetting info
		Targets:   def.TargetAll,
		DamageSrc: t.index,
		SelfHarm:  selfHarm,

		//no stats, should never get used
		Stats: make([]float64, def.EndStatType),
	}

	var mult float64
	switch typ {
	case def.Overload:
		mult = 2
		ds.Element = def.Pyro
		ds.AttackTag = def.AttackTagOverloadDamage
		ds.ICDTag = def.ICDTagOverloadDamage
		ds.ICDGroup = def.ICDGroupReactionB
		ds.StrikeType = def.StrikeTypeBlunt
	case def.Superconduct:
		mult = 0.5
		ds.Element = def.Cryo
		ds.AttackTag = def.AttackTagSuperconductDamage
		ds.ICDTag = def.ICDTagSuperconductDamage
	case def.ElectroCharged:
		mult = 1.2
		ds.Element = def.Electro
		ds.AttackTag = def.AttackTagECDamage
		ds.ICDTag = def.ICDTagECDamage
		ds.ICDGroup = def.ICDGroupReactionB
	case def.Shatter:
		mult = 1.5
		ds.Element = def.Physical
	case def.SwirlElectro:
		mult = 0.6
		ds.Element = def.Electro
		ds.AttackTag = def.AttackTagSwirlElectro
		ds.ICDTag = def.ICDTagSwirlElectro
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			ds.Targets = def.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	case def.SwirlCryo:
		mult = 0.6
		ds.Element = def.Cryo
		ds.AttackTag = def.AttackTagSwirlCryo
		ds.ICDTag = def.ICDTagSwirlCryo
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			ds.Targets = def.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	case def.SwirlHydro:
		mult = 0.6
		ds.Element = def.Hydro
		ds.AttackTag = def.AttackTagSwirlHydro
		ds.ICDTag = def.ICDTagSwirlHydro
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			ds.Targets = def.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	case def.SwirlPyro:
		mult = 0.6
		ds.Element = def.Pyro
		ds.AttackTag = def.AttackTagSwirlPyro
		ds.ICDTag = def.ICDTagSwirlPyro
		ds.Targets = t.index
		//calculate swirl'd aura durability
		if !selfHarm {
			ds.Targets = def.TargetAll
			x := res
			if res > 0.5*in.Durability {
				x = in.Durability
			}
			ds.Durability = 1.25*(x-1) + 25
		}
	default:
		//either not implemented or no dmg
		return def.Snapshot{}
	}

	//grab live EM
	char, _ := t.sim.CharByPos(in.ActorIndex)
	ds.Stats[def.EM] = char.Stat(def.EM)
	ds.Mult = mult

	return ds
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
