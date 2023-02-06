package simulation

import (
	"fmt"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/model"
)

func (sim *Simulation) CharacterDetails() []*model.Character {
	out := make([]*model.Character, len(sim.C.Player.Chars()))

	for i, v := range sim.cfg.Characters {
		m := make(map[string]int32)
		for k, v := range v.Sets {
			m[k.String()] = int32(v)
		}

		char := &model.Character{
			Name:     v.Base.Key.String(),
			Element:  v.Base.Element.String(),
			Level:    int32(v.Base.Level),
			MaxLevel: int32(v.Base.MaxLevel),
			Cons:     int32(v.Base.Cons),
			Weapon: &model.Weapon{
				Name:     v.Weapon.Key.String(),
				Refine:   int32(v.Weapon.Refine),
				Level:    int32(v.Weapon.Level),
				MaxLevel: int32(v.Weapon.MaxLevel),
			},
			Talents: &model.CharacterTalents{
				Attack: int32(v.Talents.Attack),
				Skill:  int32(v.Talents.Skill),
				Burst:  int32(v.Talents.Burst),
			},
			Sets:  m,
			Stats: v.Stats,
		}
		out[i] = char
	}

	//grab a snapshot for each char
	for i, c := range sim.C.Player.Chars() {
		snap := c.Snapshot(&combat.AttackInfo{
			Abil:      "stats-check",
			AttackTag: combat.AttackTagNone,
		})
		//convert all atk%, def% and hp% into flat amounts by tacking on base
		snap.Stats[attributes.HP] += c.Base.HP * (1 + snap.Stats[attributes.HPP])
		snap.Stats[attributes.DEF] += c.Base.Def * (1 + snap.Stats[attributes.DEFP])
		snap.Stats[attributes.ATK] += (c.Base.Atk + c.Weapon.Atk) * (1 + snap.Stats[attributes.ATKP])
		snap.Stats[attributes.HPP] = 0
		snap.Stats[attributes.DEFP] = 0
		snap.Stats[attributes.ATKP] = 0
		if sim.C.Combat.Debug {
			evt := sim.C.Log.NewEvent(
				fmt.Sprintf("%v final stats", c.Base.Key.Pretty()),
				glog.LogCharacterEvent,
				i,
			)
			for i, v := range snap.Stats {
				if v != 0 {
					evt.Write(attributes.StatTypeString[i], strconv.FormatFloat(v, 'f', 3, 32))
				}
			}
		}
		out[i].SnapshotStats = snap.Stats[:]
	}

	return out
}
